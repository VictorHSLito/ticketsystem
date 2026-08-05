package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/database"
	"github.com/victorhslito/ticketsystem/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type createUserDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateUser(c *gin.Context) {
	var createUserDTO createUserDTO

	if err := c.ShouldBindJSON(&createUserDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user domain.User
	err := database.DBConnection.Where("email = ?", createUserDTO.Email).First(&user).Error

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists!"})
		return
	}

	user.Username = createUserDTO.Username
	user.Email = createUserDTO.Email

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUserDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error trying to encrypt password, err: ", err.Error())
	}

	user.Password = string(hashedPassword)

	database.DBConnection.Create(&user)

	c.Header("Location", "/users/"+strconv.FormatUint(uint64(user.ID), 10))

	c.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}
