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
		users := storage.GetUsers(db)
		users = hidePasswords(users)
		ctx.JSON(200, users)
	}
}

func hidePasswords(users []*storage.User) []*storage.User {
	for _, u := range users {
		u.Password = "****"
	}
	return users
}

type createUserCommand struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func CreateUser(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		var command createUserCommand

		ctx.ShouldBindJSON(&command)
		log.Printf("New user username: '%s'", command.Name)

		if command.Name == "" {
			log.Printf("Empty username")
			ctx.JSON(400, "Empty username")
			return
		}
		if command.Password == "" {
			log.Printf("Empty password")
			ctx.JSON(400, "Empty password")
			return
		}

		user, err := storage.CreateUser(db, command.Name, command.Password)
		if err != nil {
			ctx.JSON(400, err.Error())
			return
		}
		ctx.JSON(200, user)
	}
}
