package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/controllers"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
)

func main() {
	database.CreateDatabaseConnection()

	database.DBConnection.AutoMigrate(&domain.Event{})

	router := gin.Default()

	router.POST("/events", controllers.CreateEvent)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed on start server, error: ", err)
	}
}
