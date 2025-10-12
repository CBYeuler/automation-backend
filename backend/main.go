package main

import (
	"log"

	"github.com/CBYeuler/automation-backend/backend/database"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	database.ConnectDatabase()

	router := gin.Default()

	// Define a simple health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})
	api := router.Group("/api/v1")
	{
		// TODO: Add machine routes (CRUD) here in the next step
		api.GET("/machines", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Machine list placeholder"})
		})
	}

	log.Println("Starting API Server on :8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
