package main

import (
	"golang-base/controllers"
	"golang-base/middlewares"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// init gets called before the main function
func init() {
	// Log error if .env file does not exist
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found")
	}
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1")
	{
		hello := new(controllers.HelloWorldController)
		v1.GET("/hello", hello.Default)

		user := new(controllers.UserController)
		v1.POST("/signup", user.Signup)
		v1.POST("/login", user.Login)
		v1.PUT("/reset-link", user.ResetLink)
		v1.PUT("/password-reset", user.PasswordReset)
		v1.PUT("/verify-link", user.VerifyLink)
		v1.PUT("/verify-account", user.VerifyAccount)
		v1.GET("/refresh", user.RefreshToken)

		bookmarks := v1.Group("/bookmarks")

		link := new(controllers.BookmarkController)
		bookmarks.Use(middlewares.Authenticate())
		{
			bookmarks.GET("/all", link.FetchBookmarks)
			bookmarks.POST("/create", link.CreateBookmak)
			bookmarks.DELETE("/delete", link.DeleteBookmark)
		}
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"message": "Not found"})
	})

	router.Run(":5000")
}
