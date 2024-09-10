package handlers

import (
	m "cdecode/pkg/models"
	s "cdecode/pkg/storage"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler = func(*gin.Context)

const defaultUsername = "admin"

func GetUsers(res s.DatabaseResolver) Handler {
	return func(ctx *gin.Context) {
		db := res()
		defer db.Close()
		users := db.GetUsers()
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
func CreateUser(res s.DatabaseResolver) Handler {
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

		db := res()
		defer db.Close()

		user, err := db.CreateUser(command.Name, command.Password)
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		ctx.JSON(200, user)
	}
}

func GetUserDataById(res s.DatabaseResolver) Handler {
	return func(ctx *gin.Context) {

		id, err := strconv.Atoi(ctx.Param("id"))

		if err != nil {
			ResultFromError(ctx, err)
			return
		}

		db := res()
		defer db.Close()

		user, err := db.GetUserById(id)
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		user.HidePassword()
		ctx.JSON(200, user)
	}
}

func GetUserData(res s.DatabaseResolver) Handler {
	return authenticated(func(ctx *gin.Context, auth *AuthState) {

		db := res()
		defer db.Close()

		user, err := db.GetUserById(auth.Claims.UserId())
		if err != nil {
			ResultFromError(ctx, err)
			return
		}
		// user.HidePassword()
		ctx.JSON(200, user)
	})
}
