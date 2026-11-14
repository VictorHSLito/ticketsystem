package main

import (
	"context"
	"errors"
	"log"

	"github.com/victorhslito/ticketsystem/internal/service"
)

func main() {
	app, err := NewApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	bookingConsumer := service.NewConsumer(
		app.SQSClient,
		service.ConsumerConfig{
			QueueURL:          app.QueueURL,
			VisibilityTimeout: 30,
			WaitTimeSeconds:   20,
			MaxMessages:       10,
			WorkerCount:       2,
		},
		func(ctx context.Context, msg *service.Message) error {
			if msg.Type != "booking" {
				return errors.New("unsupported message type")
			}

			bookingID, ok := msg.Payload["bookingId"].(float64)
			if !ok {
				return errors.New("invalid bookingId")
			}

			log.Printf("Processing booking %.0f\n", bookingID)

			app.ExpireBookingService.ExpireBookings(int(bookingID))

			log.Printf("Booked %.0f proccessed...\n", bookingID)

			return nil
		},
	)

	go func() {
		if err := bookingConsumer.Start(app.Ctx); err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("SQS consumer stopped: %v", err)
		}
	}()

	log.Println("[INFO] - Starting server on :8080")
	if err := app.Router.Run(":8080"); err != nil {
		log.Fatal("Failed on start server, error: ", err)
	}
}
