package controllers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
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
	NumberOfSeats int `json:"numberOfSeats" binding:"required"`
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

func CreateEvent(c *gin.Context) {
	var eventDTO CreateEventDTO

	if !BindAndValidate(c, &eventDTO) {
		return
	}

	event := createEvent(eventDTO)

	if event == nil {
		c.Error(ErrInternalServer)
		return
	}

	c.Header("Location", "/events/"+strconv.FormatUint(uint64(event.ID), 10))

	c.JSON(http.StatusCreated, gin.H{"message": "event created!"})
}

func createEvent(eventDTO CreateEventDTO) *domain.Event {
	var event domain.Event

	event.Name = eventDTO.Name
	event.Description = eventDTO.Description

	if err := database.DBConnection.Create(&event).Error; err != nil {
		log.Println("[ERROR] - Error trying to create event!")
		return nil
	}

	seats := createSeats(event, eventDTO.EventSettings.NumberOfSeats)
	event.Seats = seats

	return &event
}

func createSeats(event domain.Event, numberOfSeats int) []domain.Seat {
	seats := make([]domain.Seat, numberOfSeats)
	for i := range numberOfSeats {
		seats[i] = domain.Seat{
			EventID:    event.ID,
			SeatCode:   "SA" + strconv.Itoa(i),
			SeatStatus: domain.StatusAvailable,
		}
	}

	if err := database.DBConnection.Create(&seats).Error; err != nil {
		log.Println("[ERROR] - Error trying to create seats!")
		return nil
	}

	return seats
}

func ListEvents(c *gin.Context) {
	eventDTO := make([]EventsDTO, 0)
	var events []domain.Event

	database.DBConnection.Scopes(database.Paginate(c)).Find(&events)

	for _, event := range events {
		eventDTO = append(eventDTO, EventsDTO{event.Name, event.Description})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": eventDTO, "pagination": Pagination{
			c.GetInt("page"),
			c.GetInt("pageSize"),
			len(eventDTO),
			c.GetInt("totalPages"),
		},
	})
}

func ListEventByID(c *gin.Context) {
	id := c.Param("id")

	event, err := gorm.G[domain.Event](database.DBConnection).Where("id = ?", id).First(c.Request.Context())
	if err != nil {
		c.Error(ErrEventNotFound)
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": EventsDTO{event.Name, event.Description}})
}

func ListAllSeatsAvailables(c *gin.Context) {
	eventID := c.Param("id")
	var seats []domain.Seat

	query := database.DBConnection.Where("event_id = ? and seat_status = ?", eventID, domain.StatusAvailable)

	err := query.Scopes(database.Paginate(c)).Find(&seats).Error
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
