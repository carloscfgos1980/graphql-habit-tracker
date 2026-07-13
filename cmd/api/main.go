package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Habit Streak Tracker is running",
			"status":  "success",
		})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Default port if not specified in .env
	}
	router.Run(":" + port) // Start the server on the specified port
	// router.Run() // Start the server on the default port 3000
}
