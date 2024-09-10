package main

import (
	h "cdecode/pkg/handlers"
	s "cdecode/pkg/storage"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

const defaultUsername = "admin"

var db *sql.DB

func main() {
	// db setup

	db := s.GetDB()
	db.Init()

	defer db.Close()

	// server setup
	r := gin.Default()
	log.SetOutput(gin.DefaultWriter)

	// routes setup
	aliases := r.Group("/alias")
	{
		aliases.GET("/", h.GetAliases(s.GetDB))               // Get user's aliases
		aliases.POST("/new", h.CreateAlias(s.GetDB))          // Create new alias for specified resource
		aliases.DELETE("/delete/:id", h.DeleteAlias(s.GetDB)) // Delete user's alias
	}

	users := r.Group("/user")
	{
		users.GET("/", h.GetUsers(s.GetDB))                // List users
		users.GET("/data/:id", h.GetUserDataById(s.GetDB)) // Get user data
		users.GET("/data/", h.GetUserData(s.GetDB))        // Get current authorized user data
		users.POST("/new", h.CreateUser(s.GetDB))          // Register
		users.POST("/change/password", dummyHandler)       // Change password TODO: implement
		users.POST("/change/name", dummyHandler)           // Change username TODO: implement
	}

	resources := r.Group("resource")
	{
		resources.GET("/", h.GetResources(s.GetDB))                // Get users resources
		resources.POST("/create", h.CreateResource(s.GetDB))       // Create new resource
		resources.DELETE("/delete/:id", h.DeleteResource(s.GetDB)) // Delete existing resource
		resources.POST("/edit", h.UpdateResource(s.GetDB))         // Edit existing resource content
	}

	r.POST("/auth", h.Authenticate(s.GetDB))

	r.GET("/read/:name", h.Read(s.GetDB))

	r.Run(":1488")
}

func dummyHandler(ctx *gin.Context) {
	log.Printf("Dummy handler")
	ctx.JSON(200, "Dummy handler")
}
