package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/aws"
	"github.com/victorhslito/ticketsystem/internal/controllers"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/routes"
)

func main() {
	ctx := context.Background()

	database.CreateDatabaseConnection()

	router := gin.Default()
	router.Use(controllers.ErrorHandler())

	database.DBConnection.AutoMigrate(&domain.Event{})
	database.DBConnection.AutoMigrate(&domain.Seat{})
	database.DBConnection.AutoMigrate(&domain.User{})
	database.DBConnection.AutoMigrate(&domain.Booking{})
	database.DBConnection.AutoMigrate(&domain.Ticket{})

	routes.RegisterBookingGroup(router)
	routes.RegisterEventsGroup(router)
	routes.RegisterUserGroup(router)

	cfg := aws.LoadFromEnv()

	sqsClient, err := aws.NewSQSClient(ctx, cfg)
	if err != nil {
		log.Fatal("Error trying to connect in SQS", err)
	}

	queueManager := aws.NewQueueManager(sqsClient)

	queueURL, _ := queueManager.GetQueueUrl(ctx, "check-booking-pending-state")
	fmt.Println("Queue Manager Connected on Queue: ", queueURL)

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed on start server, error: ", err)
	}
}
