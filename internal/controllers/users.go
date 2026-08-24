package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type UserController struct {
	s *service.UserService
}

type CreateUserDTO struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func NewUserController(s *service.UserService) *UserController {
	return &UserController{s: s}
}

func (h *UserController) CreateUser(c *gin.Context) {
	var createUserDTO CreateUserDTO

	if !BindAndValidate(c, &createUserDTO) {
		return
	}

	userID, err := h.s.CreateUser(createUserDTO.Username, createUserDTO.Email, createUserDTO.Password)
	if err != nil {
		c.Error(ErrUserAlreadyExists)
		return
	}

	c.Header("Location", "/users/"+strconv.FormatUint(uint64(userID), 10))

	c.JSON(http.StatusCreated, gin.H{"message": "User created!"})
}
