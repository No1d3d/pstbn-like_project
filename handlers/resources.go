package handlers

import (
	"cdecode/storage"
	"database/sql"
	"net/http"

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
			ctx.JSON(http.StatusBadRequest, BadResult("Validation error"))
			return
		}

		r, err := storage.CreateResource(db, auth.Claims.UserId(), data.Content)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResultFromError(err))
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
