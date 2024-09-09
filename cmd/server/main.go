package main

import (
	"cdecode/handlers"
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
	log.SetOutput(gin.DefaultErrorWriter)

	// routes setup
	aliases := r.Group("/alias")
	{
		aliases.GET("/", getAliasesHandler)
		aliases.POST("/new", createAlias)
	}
	users := r.Group("/user")
	{
		users.GET("/", h.GetUsers(db))
		users.GET("/data/:id", h.GetUserById(db))
		users.GET("/data/", h.GetUserData(db))
		users.POST("/new", h.CreateUser(db))
	}

	r.POST("/auth", h.Authenticate(db))

	resources := r.Group("resource")
	{
		resources.GET("/", h.GetResources(db))
		resources.POST("/create", h.CreateResource(db))
	}

	r.Run(":1488")
}

func getAliasesHandler(ctx *gin.Context) {
	ctx.JSON(200, handlers.GetAliases(db))
}
func createAlias(ctx *gin.Context) {
	log.Fatalf("TODO: Implement create alias method")
}
