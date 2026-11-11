package controllers

import (
	"log/slog"
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
	logger      *slog.Logger
}

func NewAuthController(authService *service.AuthService, logger *slog.Logger) *AuthController {
	return &AuthController{authService: authService, logger: logger.With("component", "auth_handler")}
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
