package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/routes"
)

func main() {
	database.CreateDatabaseConnection()

	router := gin.Default()

	database.DBConnection.AutoMigrate(&domain.Event{})
	database.DBConnection.AutoMigrate(&domain.Seat{})
	database.DBConnection.AutoMigrate(&domain.User{})
	database.DBConnection.AutoMigrate(&domain.Booking{})
	database.DBConnection.AutoMigrate(&domain.Ticket{})

	routes.RegisterBookingGroup(router)
	routes.RegisterEventsGroup(router)
	routes.RegisterUserGroup(router)

	err := router.Run(":8080")
	if err != nil {
		log.Fatal("Failed on start server, error: ", err)
	}
}
