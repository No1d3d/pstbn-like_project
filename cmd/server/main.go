package main

import (
	"cdecode/storage"
	"database/sql"

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
	r.GET("/users", getUsersHandler)
	r.GET("/resources", getResourcesHandler)

	r.Run(":1488")
}

func getAliasesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetAliases(db, defaultUsername))
}
func getUsersHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetUsers(db, defaultUsername))
}
func getResourcesHandler(ctx *gin.Context) {
	ctx.JSON(200, storage.GetResources(db, defaultUsername))
}
