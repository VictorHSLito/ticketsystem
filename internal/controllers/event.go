package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
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
