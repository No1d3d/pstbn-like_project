package handlers

import (
	"cdecode/storage"
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler = func(*gin.Context)

const defaultUsername = "admin"
const passwordPlaceholder = "****"

func GetUsers(db *sql.DB) Handler {
	return func(ctx *gin.Context) {
		users := storage.GetUsers(db)
		users = hidePasswords(users)
		ctx.JSON(200, users)
	}
}

func hidePasswords(users []*storage.User) []*storage.User {
	for _, u := range users {
		u.Password = passwordPlaceholder
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
			ctx.JSON(400, BadResult("Empty username"))
			return
		}
		if command.Password == "" {
			log.Printf("Empty password")
			ctx.JSON(400, BadResult("Empty password"))
			return
		}

		user, err := storage.CreateUser(db, command.Name, command.Password)
		if err != nil {
			ctx.JSON(400, ResultFromError(err))
			return
		}
		ctx.JSON(200, user)
	}
}

type getUserData struct {
	Id int `url:"id"`
}

func GetUser(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		var data getUserData

		ctx.BindUri(data)

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(400, ResultFromError(err))
			return
		}
		user, err := storage.GetUserById(db, id)
		if err != nil {
			ctx.JSON(400, ResultFromError(err))
			return
		}
		user.HidePassword()
		ctx.JSON(200, user)
	}
}
