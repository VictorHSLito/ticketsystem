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
	"github.com/victorhslito/ticketsystem/internal/repository"
	"github.com/victorhslito/ticketsystem/internal/routes"
	"github.com/victorhslito/ticketsystem/internal/service"
)

func main() {
	ctx := context.Background()

	db := database.CreateDatabaseConnection()
	db.AutoMigrate(&domain.Event{})
	db.AutoMigrate(&domain.Seat{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Booking{})
	db.AutoMigrate(&domain.Ticket{})

	userRepository := repository.NewUserRepository(db)
	eventRepository := repository.NewEventRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	userService := service.NewUserService(userRepository)
	eventService := service.NewEventService(eventRepository)
	bookingService := service.NewBookingService(bookingRepository, userRepository, eventRepository)

	userHandler := controllers.NewUserController(userService)
	eventHandler := controllers.NewEventController(eventService)
	bookingHandler := controllers.NewBookingController(bookingService)

	router := gin.Default()
	router.Use(controllers.ErrorHandler())

	routes.RegisterBookingGroup(router, bookingHandler)
	routes.RegisterEventsGroup(router, eventHandler)
	routes.RegisterUserGroup(router, userHandler)

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
