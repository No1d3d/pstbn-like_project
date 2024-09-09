package handlers

import (
	"cdecode/storage"
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type createResourceData struct {
	Content string `json:"content"`
}

func (d createResourceData) Validate() bool {
	return d.Content != ""
}

func CreateResource(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		var data createResourceData
		ctx.BindJSON(&data)
		if !data.Validate() {
			BadRequest(ctx, "Validation error")
			return
		}

		r, err := storage.CreateResource(db, auth.Claims.UserId(), data.Content)

		if err != nil {
			InternalError(ctx)
			return
		}

		ctx.JSON(200, r)
	})
}

func GetResources(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		res := storage.GetResources(db, auth.Claims.UserId())
		ctx.JSON(200, res)
	})
}

func DeleteResource(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, as *AuthState) {
		id, _ := strconv.Atoi(ctx.Param("id"))

		res, err := storage.GetResourceById(db, id)
		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		log.Println("Resource found")

		err = storage.DeleteResource(db, res.Id)
		if err != nil {
			log.Println(err)
			ResultFromError(ctx, err)
			return
		}
		Success(ctx)
	})
}
