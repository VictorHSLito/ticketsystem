package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/victorhslito/ticketsystem/internal/aws"
	"github.com/victorhslito/ticketsystem/internal/controllers"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
	"github.com/victorhslito/ticketsystem/internal/routes"
	"github.com/victorhslito/ticketsystem/internal/service"
)

func main() {
	_ = godotenv.Load("../.env")

	ctx := context.Background()

	cfg := aws.LoadFromEnv()

	sqsClient, err := aws.NewSQSClient(ctx, cfg)
	if err != nil {
		log.Fatal("Error trying to connect in SQS", err)
	}

	queueManager := aws.NewQueueManager(sqsClient)

	queueURL, err := queueManager.GetQueueUrl(ctx, "check-booking-pending-state")
	if err != nil {
		log.Fatalf("failed to read Queue URL: %v", err)
	}

	fmt.Println("Queue Manager Connected on Queue: ", queueURL)

	db := database.CreateDatabaseConnection()
	db.AutoMigrate(&domain.Event{})
	db.AutoMigrate(&domain.Seat{})
	db.AutoMigrate(&domain.Role{})
	db.AutoMigrate(&domain.Scope{})
	db.AutoMigrate(&domain.User{})
	db.AutoMigrate(&domain.Booking{})
	db.AutoMigrate(&domain.Ticket{})

	userRepository := repository.NewUserRepository(db)
	eventRepository := repository.NewEventRepository(db)
	bookingRepository := repository.NewBookingRepository(db)
	seatRepository := repository.NewSeatsRepository(db)
	ticketRepository := repository.NewTicketRepository(db)

	userService := service.NewUserService(userRepository)
	eventService := service.NewEventService(eventRepository)
	bookingService := service.NewBookingService(bookingRepository, userRepository, eventRepository)
	queueService := service.NewQueueService(sqsClient, queueURL)
	confirmBookingService := service.NewConfirmBookingService(bookingRepository)
	expireBookingService := service.NewExpireBookingService(bookingRepository, seatRepository, ticketRepository)
	authService := service.NewAuthService(userRepository)

	userHandler := controllers.NewUserController(userService)
	eventHandler := controllers.NewEventController(eventService)
	bookingHandler := controllers.NewBookingController(bookingService, queueService, confirmBookingService)
	authHandler := controllers.NewAuthController(authService)

	router := gin.Default()
	router.Use(controllers.ErrorHandler())

	routes.RegisterBookingGroup(router, bookingHandler)
	routes.RegisterEventsGroup(router, eventHandler)
	routes.RegisterUserGroup(router, userHandler)
	routes.RegisterAuthGroup(router, authHandler)

	bookingConsumer := service.NewConsumer(sqsClient,
		service.ConsumerConfig{
			QueueURL:          queueURL,
			VisibilityTimeout: 30,
			WaitTimeSeconds:   20,
			MaxMessages:       10,
			WorkerCount:       2,
		},
		func(ctx context.Context, msg *service.Message) error {
			if msg.Type != "booking" {
				return fmt.Errorf("unsupported message type: %s", msg.Type)
			}

			bookingID, ok := msg.Payload["bookingId"].(float64)
			if !ok {
				return fmt.Errorf("invalid bookingId")
			}

			log.Printf("Processing booking %.0f\n", bookingID)

			expireBookingService.ExpireBookings(int(bookingID))

			log.Printf("Booked %.0f proccessed...\n", bookingID)

			return nil
		})

	go func() {
		if err := bookingConsumer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("SQS consumer stopped: %v", err)
		}
	}()

	err = router.Run(":8080")
	if err != nil {
		log.Fatal("Failed on start server, error: ", err)
	}
}
