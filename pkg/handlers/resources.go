package handlers

import (
	m "cdecode/pkg/models"
	s "cdecode/pkg/storage"
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

func CreateResource(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		var data createResourceData
		ctx.BindJSON(&data)
		if !data.Validate() {
			BadRequest(ctx, "Validation error")
			return
		}

		db := res()
		defer db.Close()

		r, err := db.CreateResource(auth.Claims.UserId(), data.Content)

		if err != nil {
			InternalError(ctx)
			return
		}

		ctx.JSON(200, r)
	})
}

func GetResources(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {
		db := res()
		defer db.Close()
		result := db.GetResources(auth.Claims.UserId())
		ctx.JSON(200, result)
	})
}

func DeleteResource(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, as *AuthState) {
		id, _ := strconv.Atoi(ctx.Param("id"))

		db := res()
		defer db.Close()

		result, err := db.GetResourceById(id)
		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}

		if result.UserId != as.Claims.UserId() {
			NotFound(ctx)
			return
		}

		err = db.DeleteResource(result.Id)
		if err != nil {
			log.Println(err)
			ResultFromError(ctx, err)
			return
		}
		Success(ctx)
	})
}

type updateResourceData struct {
	Id      m.ResourceId `json:"id"`
	Content string       `json:"content"`
}

func UpdateResource(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, as *AuthState) {
		var data updateResourceData
		ctx.BindJSON(&data)

		db := res()
		defer db.Close()

		result, err := db.GetResourceById(data.Id)
		if err != nil {
			log.Println(err)
			NotFound(ctx)
			return
		}
		if result.UserId != as.Claims.UserId() {
			NotFound(ctx)
			return
		}

		result.Content = data.Content

		err = db.UpdateResource(*result)
		if err != nil {
			ResultFromError(ctx, err)
		}
		ctx.JSON(200, result)
	})
}
