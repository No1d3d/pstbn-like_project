package handlers

import (
	"cdecode/storage"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createResourceData struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func (d createResourceData) Validate() bool {
	return d.Name != "" && d.Content != ""
}

func CreateResource(db *sql.DB) Handler {
	return func(ctx *gin.Context) {
		var data createResourceData

		if !isAuthenticated(ctx) {
			ctx.JSON(http.StatusUnauthorized, BadResult("Not authenticated"))
			return
		}

		id, err := getAuthenticatedUserId(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, BadResult("Not authenticated"))
			return
		}

		ctx.BindJSON(&data)
		if !data.Validate() {
			ctx.JSON(http.StatusBadRequest, BadResult("Validation error"))
			return
		}

		r, err := storage.CreateResource(db, id, data.Content)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResultFromError(err))
			return
		}

		ctx.JSON(200, r)
	}
}

func GetResources(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		if !isAuthenticated(ctx) {
			ctx.JSON(http.StatusUnauthorized, "")
			return
		}
		id, _ := getAuthenticatedUserId(ctx)
		res := storage.GetResources(db, id)

		ctx.JSON(200, res)
	}
}
