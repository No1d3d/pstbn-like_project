package handlers

import (
	m "cdecode/models"
	"cdecode/storage"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type createAliasData struct {
	Name       string       `json:"content"`
	ResourceId m.ResourceId `json:"resourceId"`
}

func (d createAliasData) Validate() bool {
	return d.Name != "" && d.ResourceId >= 0
}

func CreateAliases(db *sql.DB) Handler {
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

func GetAliases(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		res := storage.GetAliases(db, auth.Claims.UserId())
		ctx.JSON(200, res)
	})
}
