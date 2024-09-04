package handlers

import (
	"cdecode/storage"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

type Handler = func(*gin.Context)

const defaultUsername = "admin"

func GetUsers(db *sql.DB) Handler {
	return func(ctx *gin.Context) {
		ctx.JSON(200, storage.GetUsers(db, defaultUsername))
	}
}

type createUserCommand struct {
	Username string `json:"username"`
}

func CreateUser(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		var command createUserCommand

		ctx.ShouldBindJSON(&command)
		log.Printf("New user username: '%s'", command.Username)

		if command.Username == "" {
			log.Printf("Empty username")
			ctx.JSON(400, "Empty username")
			return
		}

		user, err := storage.CreateUser(db, defaultUsername, command.Username)
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
		ctx.JSON(200, user)
	}
}
