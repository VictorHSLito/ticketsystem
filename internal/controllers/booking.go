package controllers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type BookingController struct {
	bookingService        *service.BookingService
	queueService          *service.QueueService
	confirmBookingService *service.ConfirmBookingService
	logger                *slog.Logger
}

type CreateBookingDTO struct {
	UserID  uint     `json:"userId" binding:"required"`
	EventID uint     `json:"eventId" binding:"required"`
	Seats   []string `json:"seats" binding:"required"`
}

type ConfirmBookingDTO struct {
	BookingID uint `json:"bookingId" binding:"required"`
}

func NewBookingController(
	bookingService *service.BookingService,
	queueService *service.QueueService,
	confirmBookingService *service.ConfirmBookingService,
	logger *slog.Logger,
) *BookingController {
	return &BookingController{bookingService: bookingService, queueService: queueService, confirmBookingService: confirmBookingService, logger: logger.With("component", "booking_handler")}
}

func (h *BookingController) CreateBooking(c *gin.Context) {
	var bookingDTO CreateBookingDTO

	if !BindAndValidate(c, &bookingDTO) {
		return
	}

	booking, err := h.bookingService.CreateBooking(bookingDTO.UserID, bookingDTO.EventID, bookingDTO.Seats)
	if err != nil {
		h.logger.Error("Failed to process booking", "error", err)
		c.Error(ErrInternalServer)
		return
	}

	_, err = h.queueService.SendMessage(c.Request.Context(), &service.Message{Payload: map[string]any{
		"bookingId": booking.ID,
	}, Type: "booking"}, 10)
	if err != nil {
		h.logger.Error("Failed to send message", "error", err)
		c.Error(ErrInternalServer)
		return
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
		h.logger.Error("Failed to convert id", "error", err)
		c.Error(ErrInternalServer)
		return
	}

	booking, err := h.bookingService.FindByID(uint(id))
	if err != nil {
		c.Error(ErrBookingNotFound)
		return
	}

	c.JSON(http.StatusOK, booking)
}

func (h *BookingController) ConfirmBooking(c *gin.Context) {
	var confirmDTO ConfirmBookingDTO

	if !BindAndValidate(c, &confirmDTO) {
		c.Error(ErrInternalServer)
		return
	}

	h.confirmBookingService.ConfirmBooking(confirmDTO.BookingID)

	c.JSON(http.StatusNoContent, gin.H{})
}
