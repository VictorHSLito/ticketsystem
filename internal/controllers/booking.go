package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
)

type CreateBookingDTO struct {
	UserID  uint     `json:"userId" binding:"required"`
	EventID uint     `json:"eventId" binding:"required"`
	Seats   []string `json:"seats" binding:"required"`
}

func CreateBooking(c *gin.Context) {
	var bookingDTO CreateBookingDTO
	var user domain.User
	var event domain.Event

	if err := c.ShouldBindJSON(&bookingDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Recupera o usuário
	database.DBConnection.Where("id = ?", bookingDTO.UserID).Find(&user)

	database.DBConnection.Preload("Seats").Where("id = ?", bookingDTO.EventID).Find(&event)

	// Cria o dado de booking
	booking := domain.Booking{
		UserID:        bookingDTO.UserID,
		BookingStatus: domain.Pending,
	}

	database.DBConnection.Create(&booking)

	// Cria os tickets e vincula a reserva
	createTickets(event.Seats, bookingDTO.Seats, booking.ID)

	// Anexa o booking com o usuário
	user.Bookings = append(user.Bookings, booking)

	c.Header("Location", "/booking/"+strconv.FormatUint(uint64(booking.ID), 10))
	c.JSON(http.StatusCreated, gin.H{
		"message": "Booked created!",
	})
}

func ListBookingByID(c *gin.Context) {
	bookingID := c.Param("id")

	var booking domain.Booking
	database.DBConnection.Where("id = ?", bookingID).Find(&booking)
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
		seat.SeatStatus = domain.StatusBooked
	}

	return bookingSeats
}
