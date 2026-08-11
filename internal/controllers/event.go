package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/service"
)

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

func CreateEvent(c *gin.Context, s *service.EventService) {
	var eventDTO CreateEventDTO

	if !BindAndValidate(c, &eventDTO) {
		return
	}

	eventID := s.CreateEvent(eventDTO.Name, eventDTO.Description, eventDTO.EventSettings.NumberOfSeats)

	c.Header("Location", "/events/"+strconv.FormatUint(uint64(eventID), 10))

	c.JSON(http.StatusCreated, gin.H{"message": "event created!"})
}

func ListEvents(c *gin.Context, s *service.EventService) {
	var eventDTO []EventsDTO
	events := s.ListEvents()

	for _, event := range events {
		eventDTO = append(eventDTO, EventsDTO{event.Name, event.Description})
	}

	database.Paginate(c)
	c.JSON(http.StatusOK, gin.H{
		"data": eventDTO, "pagination": Pagination{
			c.GetInt("page"),
			c.GetInt("pageSize"),
			len(eventDTO),
			c.GetInt("totalPages"),
		},
	})
}

func ListEventByID(c *gin.Context, s *service.EventService) {
	id := c.Param("id")
	eventID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		fmt.Printf("could not parse id: %v\n", err)
		return
	}

	event := s.ListByID(uint(eventID))

	if event == nil {
		c.Error(ErrEventNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": EventsDTO{event.Name, event.Description}})
}

func ListAllSeatsAvailables(c *gin.Context, s *service.EventService) {
	id := c.Param("id")
	eventID, err := strconv.ParseUint(id, 10, 64)
	seats := s.ListSeatsAvailables(uint(eventID))

	database.Paginate(c)
	if err != nil {
		c.Error(ErrEventNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": seats, "pagination": Pagination{
		PageNumber: c.GetInt("page"),
		PageSize:   c.GetInt("pageSize"),
		TotalItems: len(seats),
		TotalPages: c.GetInt("totalPages"),
	}})
}
