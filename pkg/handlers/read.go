package handlers

import (
	"cdecode/pkg/storage"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

func Read(db *sql.DB) Handler {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		alias, err := storage.GetAliasByName(db, name)
		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		res, err := storage.GetResourceById(db, alias.ResourceId)

		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		ctx.JSON(200, res)
	}
}
