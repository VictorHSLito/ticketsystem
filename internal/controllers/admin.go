package controllers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/service"
)

type AdminController struct {
	adminService *service.AdminService
	logger       *slog.Logger
}

type CreateAdminRequestDTO struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

func NewAdminController(adminService *service.AdminService, logger *slog.Logger) *AdminController {
	return &AdminController{adminService: adminService, logger: logger.With("component", "admin_handler")}
}

func (h *AdminController) CreateAdmin(c *gin.Context) {
	var createAdminDTO CreateAdminRequestDTO

	if !BindAndValidate(c, &createAdminDTO) {
		return
	}

	err := h.adminService.SetupAdminUser(createAdminDTO.Username, createAdminDTO.Password, createAdminDTO.Email)
	if err != nil {
		c.Error(ErrInternalServer)
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
