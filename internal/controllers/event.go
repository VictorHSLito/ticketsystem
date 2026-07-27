package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type EventDTO struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required"`
}

func CreateEvent(c *gin.Context) {
	var eventDTO EventDTO

	if err := c.ShouldBindJSON(&eventDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	event := &domain.Event{Name: eventDTO.Name, Description: eventDTO.Description}

	if err := database.DBConnection.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "event created!"})
}

func ListEvents(c *gin.Context) {
	var eventDTO []EventDTO
	var events []domain.Event

	database.DBConnection.Find(&events)

	for _, event := range events {
		eventDTO = append(eventDTO, EventDTO{event.Name, event.Description})
	}

	c.JSON(http.StatusOK, gin.H{"data": eventDTO})
}

func ListEventByID(c *gin.Context) {
	id := c.Param("id")

	event, err := gorm.G[domain.Event](database.DBConnection).Where("id = ?", id).First(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": EventDTO{event.Name, event.Description}})
}
