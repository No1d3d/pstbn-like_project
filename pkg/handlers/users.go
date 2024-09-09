package handlers

import (
	m "cdecode/pkg/models"
	s "cdecode/pkg/storage"
	"database/sql"
	"log"
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

// Registration endpoint
func CreateUser(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		var command createUserCommand

		ctx.ShouldBindJSON(&command)
		log.Printf("New user username: '%s'", command.Name)

		if command.Name == "" {
			BadRequest(ctx, "Empty username")
			return
		}
		if command.Password == "" {
			log.Printf("Empty password")
			BadRequest(ctx, "Empty password")
			return
		}

		user, err := s.CreateUser(db, command.Name, command.Password)
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		ctx.JSON(200, user)
	}
}

func GetUserDataById(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		user, err := s.GetUserById(db, id)
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		user.HidePassword()
		ctx.JSON(200, user)
	}
}

func GetUserData(db *sql.DB) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {

		user, err := s.GetUserById(db, auth.Claims.UserId())
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		// user.HidePassword()
		ctx.JSON(200, user)
	})
}
