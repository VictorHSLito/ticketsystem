package controllers

import (
	"log/slog"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type EventController struct {
	s      *service.EventService
	logger *slog.Logger
}

type CreateEventDTO struct {
	Name          string        `json:"name" binding:"required"`
	Description   string        `json:"description" binding:"required"`
	EventSettings EventSettings `json:"eventSettings" binding:"required"`
}

type EventsDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

type EventSettings struct {
	NumberOfSeats int `json:"numberOfSeats" binding:"required,gte=10"`
}

type SeatDTO struct {
	SeatID     string            `json:"seatId"`
	SeatStatus domain.SeatStatus `json:"seatStatus"`
}

type Pagination struct {
	PageNumber int `json:"pageNumber"`
	PageSize   int `json:"pageSize"`
	TotalItems int `json:"totalItems"`
	TotalPages int `json:"totalPages"`
}

func NewEventController(s *service.EventService, logger *slog.Logger) *EventController {
	return &EventController{s: s, logger: logger.With("component", "event_controller")}
}

func (h *EventController) CreateEvent(c *gin.Context) {
	var eventDTO CreateEventDTO

	if !BindAndValidate(c, &eventDTO) {
		return
	}

	eventID := h.s.CreateEvent(eventDTO.Name, eventDTO.Description, eventDTO.EventSettings.NumberOfSeats)

	c.Header("Location", "/events/"+strconv.FormatUint(uint64(eventID), 10))

	c.JSON(http.StatusCreated, gin.H{"message": "Event created!"})
}

func (h *EventController) ListEvents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	events, totalElements := h.s.ListEvents(page, pageSize)

	totalPages := int(math.Ceil(float64(totalElements) / float64(pageSize)))

	var eventDTO []EventsDTO
	for _, event := range events {
		eventDTO = append(eventDTO, EventsDTO{event.Name, event.Description})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": eventDTO,
		"pagination": Pagination{
			PageNumber: page,
			PageSize:   pageSize,
			TotalItems: len(eventDTO),
			TotalPages: totalPages,
		},
	})
}

func (h *EventController) ListEventByID(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error("Could not parse eventID to uint", "error", err)
		c.Error(ErrInternalServer)
		return
	}
	event := h.s.ListByID(uint(eventID))

	if event == nil {
		c.Error(ErrEventNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": EventsDTO{event.Name, event.Description}})
}

func (h *EventController) ListAllSeats(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error("Could not parse eventID to uint", "error", err)
		c.Error(ErrInternalServer)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	seats, totalElements := h.s.ListAllSeats(uint(eventID), page, pageSize)

	totalPages := int(math.Ceil(float64(totalElements) / float64(pageSize)))
	c.JSON(http.StatusOK, gin.H{"data": seats, "pagination": Pagination{
		PageNumber: page,
		PageSize:   pageSize,
		TotalItems: len(seats),
		TotalPages: totalPages,
	}})
}

func (h *EventController) ListSeatsByStatus(c *gin.Context) {
	id := c.Param("id")
	eventID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		h.logger.Error("Could not parse eventID to uint", "error", err)
		c.Error(ErrInternalServer)
		return
	}

	status := c.DefaultQuery("seatStatus", string(domain.StatusAvailable))

	seats := h.s.ListSeatsByStatus(uint(eventID), status)

	c.JSON(http.StatusOK, gin.H{
		"data": seats,
	})
}
