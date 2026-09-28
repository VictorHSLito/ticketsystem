package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type LoginRequestDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthController struct {
	authService *service.AuthService
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (h *AuthController) GenerateToken(c *gin.Context) {
	var loginRequest LoginRequestDTO

	if !BindAndValidate(c, &loginRequest) {
		return
	}

	token, err := h.authService.GenerateToken(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.Error(ErrUserNotFound)
		return
	}

	c.JSON(http.StatusOK, token)
}
