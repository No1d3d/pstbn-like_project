package handlers

import (
	m "cdecode/models"
	s "cdecode/storage"
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler = func(*gin.Context)

const defaultUsername = "admin"

func GetUsers(db *sql.DB) Handler {
	return func(ctx *gin.Context) {
		users := s.GetUsers(db)
		users = hidePasswords(users)
		ctx.JSON(200, users)
	}
}

func hidePasswords(users []*m.User) []*m.User {
	for _, u := range users {
		u.HidePassword()
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

		user, err := s.CreateUser(db, command.Name, command.Password)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, ResultFromError(err))
			return
		}
		ctx.JSON(200, user)
	}
}

func GetUserById(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ctx.JSON(400, ResultFromError(err))
			return
		}
		user, err := s.GetUserById(db, id)
		if err != nil {
			ctx.JSON(400, ResultFromError(err))
			return
		}
		user.HidePassword()
		ctx.JSON(200, user)
	}
}

func GetUserData(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		if !isAuthenticated(ctx) {
			ctx.JSON(http.StatusBadRequest, BadResult("Not authorized"))
			return
		}

		id, err := getAuthenticatedUserId(ctx)

		if err != nil || id <= 0 {
			ctx.JSON(http.StatusInternalServerError, BadResult("Something wrong with cookie"))
			return
		}

		user, err := s.GetUserById(db, id)
		if err != nil {
			ctx.JSON(400, ResultFromError(err))
			return
		}
		// user.HidePassword()
		ctx.JSON(200, user)
	}
}
