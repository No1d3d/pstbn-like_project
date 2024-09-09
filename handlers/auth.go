package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/golang-jwt/jwt/v5"
	"github.com/dgrijalva/jwt-go"

	"cdecode/models"
	s "cdecode/storage"
)

const (
	AuthHeader = "Authorization"
)

var (
	AuthKey = []byte("Some cool key")
)

func Authenticate(db *sql.DB) Handler {
	return func(ctx *gin.Context) {

		var data AuthData
		ctx.BindJSON(&data)

		notAuthResponse := BadResult("Wrong password or username")

		user, err := s.GetUserByName(db, data.Name)
		if err != nil {
			ctx.JSON(400, notAuthResponse)
			log.Println(err)
			return
		}

		if !user.ValidatePasssword(data.Password) {
			ctx.JSON(400, notAuthResponse)
			log.Println("Not valid password")
			return
		}

		token, err := createToken(user)
		if err != nil {
			ctx.JSON(400, notAuthResponse)
			log.Println(err)
			log.Println("Something wrong with token")
			return
		}
		// TODO: add more complex logic for hashing and storing other neccessaru data in cookies
		// authenticate(ctx, user)
		ctx.JSON(200, AuthResponse{
			Token:   token,
			Success: true,
		})
	}
}

func createToken(user *models.User) (string, error) {

	claims := ApplicationClaims{
		StandardClaims: &jwt.StandardClaims{
			Id: strconv.Itoa(user.Id),
			// Issuer:    "localhost:1488",
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// return token.SignedString([]byte("sdrfjmgjknmedskjgbnfd#"))
	// r, err := getJwtKey()
	// if err != nil {
	// return "", err
	// }
	return token.SignedString(AuthKey)
}

type AuthData struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type AuthState struct {
	Claims *ApplicationClaims
}

func noAuth() AuthState {
	return AuthState{}
}

func (a *AuthState) IsAuthenticated() bool {
	return a.Claims != nil
}

func getAuthState(ctx *gin.Context) AuthState {
	header := ctx.GetHeader(AuthHeader)
	if header == "" {
		log.Printf("No authorization header ('%s')", AuthHeader)
		return noAuth()
	}
	log.Printf("Auth header: '%s'", header)

	s := strings.Split(header, " ")
	if len(s) < 2 {
		log.Printf("Not valid header, less than two words in it")
		return noAuth()
	}

	log.Printf("Token string: '%s'", s[1])

	var claims ApplicationClaims
	_, err := jwt.ParseWithClaims(s[1], &claims, func(t *jwt.Token) (interface{}, error) {
		return AuthKey, nil
	})
	if err != nil {
		log.Printf("Error while parsing JWT token: %v", err)
		return noAuth()
	}

	return AuthState{
		Claims: &claims,
	}
}

type ApplicationClaims struct {
	*jwt.StandardClaims
}

func (c *ApplicationClaims) UserId() int {
	id, _ := strconv.Atoi(c.Id)
	return id
}

func getAuthenticatedUserId(ctx *gin.Context) (models.UserId, error) {
	value := ctx.GetHeader(AuthHeader)

	id, err := strconv.Atoi(value)

	if err != nil {
		return -1, err
	}

	return id, nil
}

type AuthHandler = func(*gin.Context, *AuthState)

func authenticated(h AuthHandler) Handler {
	return func(ctx *gin.Context) {
		state := getAuthState(ctx)

		if !state.IsAuthenticated() {
			ctx.JSON(http.StatusUnauthorized, "Not authorized")
			return
		}

		h(ctx, &state)
	}
}
