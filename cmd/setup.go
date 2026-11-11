package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/victorhslito/ticketsystem/internal/aws"
	"github.com/victorhslito/ticketsystem/internal/controllers"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/logger"
	"github.com/victorhslito/ticketsystem/internal/routes"

	"github.com/victorhslito/ticketsystem/internal/repository"
	"github.com/victorhslito/ticketsystem/internal/service"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type App struct {
	Router               *gin.Engine
	Ctx                  context.Context
	QueueManager         *aws.QueueManager
	QueueURL             string
	SQSClient            *sqs.Client
	ExpireBookingService *service.ExpireBookingService
}

func NewApp() (*App, error) {
	ctx := context.Background()

	if err := godotenv.Load("../.env"); err != nil {
		fmt.Println("[ERROR] - Failed trying to read \".env\" file")
		log.Printf("error: %v\n", err)
		return nil, errors.New("verify if \".env\" file exists or the path is correct")
	}

	privateKeyPath := os.Getenv("JWT_PRIVATE_KEY_PATH")
	publicKeyPath := os.Getenv("JWT_PUBLIC_KEY_PATH")
	expiresInStr := os.Getenv("JWT_EXPIRES_IN")

	expiresInSec, err := strconv.Atoi(expiresInStr)
	if err != nil {
		expiresInSec = 3600
	}
	expiresIn := time.Duration(expiresInSec) * time.Second

	signBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, errors.New("error trying to read private key")
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return nil, errors.New("error trying to convert private key")
	}

	verifyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, errors.New("error trying to read public key")
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return nil, errors.New("error trying to convert public key")
	}

	queueName := os.Getenv("AWS_SQS_NAME")

	cfg := aws.LoadFromEnv()
	sqsClient, err := aws.NewSQSClient(ctx, cfg)
	if err != nil {
		log.Printf("[ERROR] - Error trying to connect in AWS SQS: %v", err)
		return nil, errors.New("verify if credentials are correct")
	}

	queueManager := aws.NewQueueManager(sqsClient)
	queueURL, err := queueManager.GetQueueUrl(ctx, queueName)
	if err != nil {
		log.Printf("[ERROR] - Failed to read Queue URL: %v", err)
		return nil, errors.New("verify queue credentials")
	}
	log.Println("[INFO] - Queue Manager Connected on Queue URL: ", queueURL)

	db := database.CreateDatabaseConnection()
	if err := db.AutoMigrate(
		&domain.Event{},
		&domain.Seat{},
		&domain.Role{},
		&domain.Scope{},
		&domain.User{},
		&domain.Booking{},
		&domain.Ticket{},
	); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate database: %w", err)
	}

	err = database.SeedDatabase(db)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate seed: %w", err)
	}

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
	authService := service.NewAuthService(userRepository, signKey, verifyKey, expiresIn)
	adminService := service.NewAdminService(userRepository)

	apiLogger := logger.NewLogger()

	userHandler := controllers.NewUserController(userService, apiLogger)
	eventHandler := controllers.NewEventController(eventService, apiLogger)
	bookingHandler := controllers.NewBookingController(bookingService, queueService, confirmBookingService, apiLogger)
	authHandler := controllers.NewAuthController(authService, apiLogger)
	adminHandler := controllers.NewAdminController(adminService, apiLogger)

	router := gin.Default()
	router.Use(controllers.ErrorHandler())

	routes.RegisterBookingGroup(router, bookingHandler)
	routes.RegisterEventsGroup(router, eventHandler)
	routes.RegisterUserGroup(router, userHandler)
	routes.RegisterAuthGroup(router, authHandler)
	routes.RegisterAdminGroup(router, adminHandler)

	return &App{
		Router:               router,
		Ctx:                  ctx,
		QueueManager:         queueManager,
		QueueURL:             queueURL,
		SQSClient:            sqsClient,
		ExpireBookingService: expireBookingService,
	}, nil
}
