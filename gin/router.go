package main

import (
	"frontrx/handlers"
	"frontrx/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitializeRouter initialises all the routes in the app
func InitializeRouter() {

	// ----- PUBLIC APIs -----
	// everyone can access them

	api := router.Group("/api")

	// info message
	api.GET("/", handlers.Index)

	// register users
	api.POST("/user/register", handlers.AddUser)

	// ----- PRIVATE APIs -----

	// Users
	user := api.Group("/user")
	user.Use(security.AuthJWTUser().MiddlewareFunc())
	{
		user.POST("/upload_kyc", handlers.UploadKYCDocument)
		user.GET("/me", handlers.GetUserByToken)

	}

	// ----- OTHER -----

	// In case no route is found
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "api endpoint not found"})
	})

}
