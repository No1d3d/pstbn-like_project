package handlers

import (
	m "cdecode/pkg/models"
	"cdecode/pkg/storage"
	"database/sql"
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

func CreateAlias(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		var data createAliasData
		ctx.BindJSON(&data)
		if !data.Validate() {
			BadRequest(ctx, "Validation error")
			return
		}

		r, err := storage.CreateAlias(db, auth.Claims.UserId(), data.Name, data.ResourceId)

		if err != nil {
			ResultFromError(ctx, err)
			return
		}

		ctx.JSON(200, r)
	})
}

func DeleteAlias(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {

		id, _ := strconv.Atoi(ctx.Param("id"))

		alias, err := storage.GetAliasById(db, id)
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

		err = storage.DeleteAlias(db, id)
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		Success(ctx)

	})
}

func GetAliases(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		res := storage.GetAliases(db, auth.Claims.UserId())
		ctx.JSON(200, res)
	})
}
