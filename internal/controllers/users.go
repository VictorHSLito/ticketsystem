package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type CreateUserDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(c *gin.Context, s *service.UserService) {
	var createUserDTO CreateUserDTO

	if !BindAndValidate(c, &createUserDTO) {
		return
	}

	userID := s.CreateUser(createUserDTO.Username, createUserDTO.Email, createUserDTO.Password)

	c.Header("Location", "/users/"+strconv.FormatUint(uint64(userID), 10))

	c.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}
