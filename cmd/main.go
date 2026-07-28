package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/controllers"
	"github.com/victorhslito/ticketsystem/internal/database"
)

func main() {
	database.CreateDatabaseConnection()

	router := gin.Default()

	router.GET("/events", controllers.ListEvents)
	router.GET("/events/:id", controllers.ListEventByID)
	router.POST("/events", controllers.CreateEvent)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed on start server, error: ", err)
	}
}
