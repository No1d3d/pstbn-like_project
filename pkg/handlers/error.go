package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Result struct {
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

func Success(c *gin.Context) {
	r := Result{
		Success: true,
	}
	c.JSON(http.StatusOK, r)
}

func ResultFromError(c *gin.Context, err error) {
	r := Result{
		Success:      false,
		ErrorMessage: err.Error(),
	}
	c.JSON(http.StatusBadRequest, r)
}

func BadRequest(c *gin.Context, message string) {
	r := Result{
		Success:      false,
		ErrorMessage: message,
	}

	c.JSON(http.StatusBadRequest, r)
}
func NotAuthorized(c *gin.Context) {
	r := Result{
		Success:      false,
		ErrorMessage: "Not authorized",
	}
	c.JSON(http.StatusUnauthorized, r)
}
func LoginError(c *gin.Context) {
	r := Result{
		Success:      false,
		ErrorMessage: "Can't authorize, wrong password or username",
	}
	c.JSON(http.StatusNotFound, r)
}
func NotFound(c *gin.Context) {
	r := Result{
		Success:      false,
		ErrorMessage: "Not found",
	}
	c.JSON(http.StatusNotFound, r)
}
func InternalError(c *gin.Context) {
	r := Result{
		Success:      false,
		ErrorMessage: "Internal server error",
	}
	c.JSON(http.StatusInternalServerError, r)
}
