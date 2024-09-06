package main

import (
	h "cdecode/handlers"
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

	r.GET("/users", h.GetUsers(db))
	r.GET("/user/data/:id", h.GetUserById(db))
	r.GET("/user/data/", h.GetUserData(db))
	r.POST("/user/new", h.CreateUser(db))

	r.POST("/auth", h.Auth(db))

	r.GET("/resources", h.GetResources(db))
	r.POST("/resource/create", h.CreateResource(db))

	r.Run(":1488")
}

func getAliasesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetAliases(db, defaultUsername))
}
func createAlias(ctx *gin.Context) {
	log.Fatalf("TODO: Implement create alias method")
}

// func getResourcesHandler(ctx *gin.Context) {
// ctx.JSON(200, storage.GetResources(db, defaultUsername))
// }
