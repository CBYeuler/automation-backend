package main

import (
	"log"

	"github.com/CBYeuler/automation-backend/backend/database"
	"github.com/CBYeuler/automation-backend/backend/handler"
	"github.com/CBYeuler/automation-backend/backend/repository"
	"github.com/CBYeuler/automation-backend/backend/service"
	"github.com/CBYeuler/automation-backend/backend/simulation"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize the database connection
	database.ConnectDatabase()

	db := database.GetDB()

	machineRepo := repository.NewMachineRepository(db)
	machineService := service.NewMachineService(machineRepo)
	machineHandler := handler.NewMachineHandler(machineService)

	machineSimulator := simulation.NewMachineSimulator(machineRepo)
	machineSimulator.StartGlobalSimulation()

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

		api.POST("/machines", machineHandler.CreateMachine)
		api.GET("/machines", machineHandler.GetMachines)
		api.GET("/machines/:id", machineHandler.GetMachineByID)
		api.PUT("/machines/:id", machineHandler.UpdateMachine)
		api.DELETE("/machines/:id", machineHandler.DeleteMachine)

		// Placeholder route to verify server is running
		// api.GET("/machines", func(c *gin.Context) {
		//	 c.JSON(200, gin.H{"message": "Machine list placeholder"})
		// })
	}

	log.Println("Starting API Server on :8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
