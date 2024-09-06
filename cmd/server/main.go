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
	r.GET("/user/:id", h.GetUser(db))
	r.POST("/user", h.CreateUser(db))
	r.GET("/resources", getResourcesHandler)

	r.Run(":1488")
}

func getAliasesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetAliases(db, defaultUsername))
}
func createAlias(ctx *gin.Context) {
	log.Fatalf("TODO: Implement create alias method")
}
func getResourcesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetResources(db, defaultUsername))
}
