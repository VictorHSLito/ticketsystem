package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type CreateBookingDTO struct {
	UserID  uint     `json:"userId" binding:"required"`
	EventID uint     `json:"eventId" binding:"required"`
	Seats   []string `json:"seats" binding:"required"`
}

func CreateBooking(c *gin.Context) {
	var bookingDTO CreateBookingDTO

	if err := c.ShouldBindJSON(&bookingDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Recupera o usuário
	var user domain.User
	if err := database.DBConnection.First(&user, bookingDTO.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found!"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	// Recupera o evento
	var event domain.Event
	if err := database.DBConnection.Preload("Seats").First(&event, bookingDTO.EventID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Event not found!"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
		return
	}

	bookingSeats := verifySeats(event.Seats, bookingDTO.Seats)

	if len(bookingSeats) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid seats or available seats"})
		return
	}

	// Cria o dado de booking
	booking, err := createBookingData(bookingDTO.UserID, bookingSeats)
	if err != nil {
		log.Println("[ERROR] - Failed to process booking transaction:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not process booking"})
		return
	}

	// Anexa o booking com o usuário
	user.Bookings = append(user.Bookings, booking)

	c.Header("Location", "/booking/"+strconv.FormatUint(uint64(booking.ID), 10))
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Booking created successfully!",
		"bookingId": booking.ID,
	})
}

func ListBookingByID(c *gin.Context) {
	bookingID := c.Param("id")

	var booking domain.Booking
	database.DBConnection.Where("id = ?", bookingID).Find(&booking)
}

func createBookingData(userID uint, bookedSeats map[uint]string) (domain.Booking, error) {
	var booking domain.Booking

	database.DBConnection.Transaction(func(tx *gorm.DB) error {
		booking = domain.Booking{
			UserID:        userID,
			BookingStatus: domain.Pending,
		}

		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		tickets := make([]domain.Ticket, 0, len(bookedSeats))
		var seatIDs []uint
		for seatID := range bookedSeats {
			seatIDs = append(seatIDs, seatID)
			tickets = append(tickets, domain.Ticket{
				ExternalID: uuid.NewString(),
				BookingID:  booking.ID,
				SeatID:     seatID,
			})
		}

		if err := tx.Create(&tickets).Error; err != nil {
			return err
		}

		// Atualiza o status dos assentos
		if err := tx.Model(&domain.Seat{}).Where("id IN ?", seatIDs).Update("seat_status", domain.StatusBooked).Error; err != nil {
			return err
		}

		return nil
	})

	return booking, nil
}

func createTickets(seats []domain.Seat, seatsCodes []string, bookingID uint) {
	bookedSeats := verifySeats(seats, seatsCodes)
	tickets := make([]domain.Ticket, 0, len(bookedSeats))

	for seatID := range bookedSeats {
		tickets = append(tickets, domain.Ticket{
			ExternalID: uuid.NewString(),
			BookingID:  bookingID,
			SeatID:     seatID,
		})
	}

	if len(tickets) > 0 {
		if err := database.DBConnection.Create(&tickets).Error; err != nil {
			log.Println("[ERROR] - Error trying to create tickets:", err)
			return
		}

		var seatIDs []uint
		for seatID := range bookedSeats {
			seatIDs = append(seatIDs, seatID)
		}

		if err := database.DBConnection.Model(&domain.Seat{}).Where("id IN ?", seatIDs).Update("seat_status", domain.StatusBooked).Error; err != nil {
			log.Println("[ERROR] - Error trying to update seat status:", err)
		}
	}
}

func verifySeats(seats []domain.Seat, seatCodes []string) map[uint]string {
	bookingSeats := make(map[uint]string)

	for _, seat := range seats {
		for _, seatCode := range seatCodes {
			if seat.SeatCode == seatCode {
				bookingSeats[seat.ID] = seatCode
			}
		}
	}

	return bookingSeats
}
