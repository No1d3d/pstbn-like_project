package main

import (
	"cdecode/storage"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

const defaultUsername = "admin"

var db *sql.DB

func main() {
	// db setup
	db = storage.InitDB()
	defer db.Close()

	// server setup
	r := gin.Default()

	// routes setup
	r.GET("/aliases", getAliasesHandler)
	r.POST("/alias", createAlias)
	r.GET("/users", getUsersHandler)
	r.POST("/user", createUser)
	r.GET("/resources", getResourcesHandler)

	r.Run(":1488")
}

func getAliasesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetAliases(db, defaultUsername))
}
func createAlias(ctx *gin.Context) {
	log.Fatalf("TODO: Implement create alias method")
}
func getUsersHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetUsers(db, defaultUsername))
}

type createUserCommand struct {
	Username string `json:"username"`
}

func createUser(ctx *gin.Context) {

	var command createUserCommand

	ctx.ShouldBindJSON(&command)
	log.Printf("New user username: '%s'", command.Username)

	if command.Username == "" {
		log.Printf("Empty username")
		ctx.JSON(400, "Empty username")
		return
	}

	storage.CreateUser(db, defaultUsername, command.Username)
	ctx.JSON(200, true)
}
func getResourcesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetResources(db, defaultUsername))
}
