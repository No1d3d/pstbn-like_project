package handlers

import (
	m "cdecode/pkg/models"
	s "cdecode/pkg/storage"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createAliasData struct {
	Name       string       `json:"name"`
	ResourceId m.ResourceId `json:"resourceId"`
}

func (d createAliasData) Validate() bool {
	return d.Name != "" && d.ResourceId >= 0
}

func CreateAlias(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		var data createAliasData
		ctx.BindJSON(&data)
		if !data.Validate() {
			BadRequest(ctx, "Validation error")
			return
		}

		db := res()
		defer db.Close()

		r, err := db.CreateAlias(auth.Claims.UserId(), data.Name, data.ResourceId)

		if err != nil {
			ResultFromError(ctx, err)
			return
		}

		ctx.JSON(200, r)
	})
}

func DeleteAlias(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {

		id, _ := strconv.Atoi(ctx.Param("id"))

		db := res()
		defer db.Close()

		alias, err := db.GetAliasById(id)
		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		if alias.CreatorId != auth.Claims.UserId() {
			log.Println("User tried to delete alias that he is not own")
			NotFound(ctx)
			return
		}

		err = db.DeleteAlias(id)
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		Success(ctx)

	})
}

func GetAliases(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		db := res()
		defer db.Close()
		result := db.GetAliases(auth.Claims.UserId())
		ctx.JSON(200, result)
	})
}
