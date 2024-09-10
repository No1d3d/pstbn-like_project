package handlers

import (
	s "cdecode/pkg/storage"
	"log"

	"github.com/gin-gonic/gin"
)

func Read(res s.DatabaseResolver) Handler {
	return func(ctx *gin.Context) {
		name := ctx.Param("name")
		db := res()
		defer db.Close()
		alias, err := db.GetAliasByName(name)
		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		res, err := db.GetResourceById(alias.ResourceId)

		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		ctx.JSON(200, res)
	}
}
