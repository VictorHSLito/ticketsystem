package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type BookingController struct {
	bookingService *service.BookingService
	queueService   *service.QueueService
}

type CreateBookingDTO struct {
	UserID  uint     `json:"userId" binding:"required"`
	EventID uint     `json:"eventId" binding:"required"`
	Seats   []string `json:"seats" binding:"required"`
}

func NewBookingController(bookingService *service.BookingService, queueService *service.QueueService) *BookingController {
	return &BookingController{bookingService: bookingService, queueService: queueService}
}

func (h *BookingController) CreateBooking(c *gin.Context) {
	var bookingDTO CreateBookingDTO

	if !BindAndValidate(c, &bookingDTO) {
		return
	}

	booking, err := h.bookingService.CreateBooking(bookingDTO.UserID, bookingDTO.EventID, bookingDTO.Seats)
	if err != nil {
		log.Println("[ERROR] - Failed to process booking:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = h.queueService.SendMessage(c.Request.Context(), &service.Message{Payload: map[string]any{
		"bookingId": booking.ID,
	}, Type: "booking"}, 10)
	if err != nil {
		log.Println("[ERROR] - Failed to send message:", err)
	}

	c.Header("Location", "/booking/"+strconv.FormatUint(uint64(booking.ID), 10))
	c.JSON(http.StatusCreated, gin.H{
		"message":   "Booking created successfully!",
		"bookingId": booking.ID,
	})
}

func (h *BookingController) ListBookingByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	booking, err := h.bookingService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	c.JSON(http.StatusOK, booking)
}
