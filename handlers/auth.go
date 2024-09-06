package handlers

import (
	"database/sql"
	"strconv"

	"github.com/gin-gonic/gin"

	"cdecode/models"
	s "cdecode/storage"
)

const (
	// AuthCookieName = "OurCoolAuthCookie"
	// CookiePath     = "/"
	// CookieDomain   = "localhost"
	AuthHeader = "Auth"
)

func Auth(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		var data AuthData
		ctx.BindJSON(&data)

		notAuthResponse := BadResult("Wrong password or username")

		user, err := s.GetUserByName(db, data.Name)
		if err != nil {
			ctx.JSON(400, notAuthResponse)
			return
		}

		if !user.ValidatePasssword(data.Password) {
			ctx.JSON(400, notAuthResponse)
			return
		}
		// TODO: add more complex logic for hashing and storing other neccessaru data in cookies
		// authenticate(ctx, user)
		ctx.JSON(200, AuthResponse{
			Token:   strconv.Itoa(user.Id),
			Success: true,
		})
	}
}

type AuthData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

func isAuthenticated(ctx *gin.Context) bool {
	value := ctx.GetHeader(AuthHeader)

	if value == "" {
		return false
	}
	// TODO: add more complex logic for cookie validation
	return true
}

func getAuthenticatedUserId(ctx *gin.Context) (models.UserId, error) {
	value := ctx.GetHeader(AuthHeader)

	id, err := strconv.Atoi(value)

	if err != nil {
		return -1, err
	}

	return id, nil
}
