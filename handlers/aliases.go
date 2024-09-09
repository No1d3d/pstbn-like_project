package handlers

import (
	"cdecode/storage"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createAliasData struct {
	Content string `json:"content"`
}

func (d createAliasData) Validate() bool {
	return d.Content != ""
}

func CreateAliases(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		var data createAliasData
		ctx.BindJSON(&data)
		if !data.Validate() {
			ctx.JSON(http.StatusBadRequest, BadResult("Validation error"))
			return
		}

		r, err := storage.CreateAlias(db, auth.Claims.UserId(), data.Content)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResultFromError(err))
			return
		}

		ctx.JSON(200, r)
	})
}

func GetAliases(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		res := storage.GetAliases(db, auth.Claims.UserId())
		ctx.JSON(200, res)
	})
}
