package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/controllers"
	"github.com/victorhslito/ticketsystem/internal/middlewares"
)

func RegisterEventsGroup(rg *gin.Engine, h *controllers.EventController) {
	events := rg.Group("/events")
	events.Use(middlewares.AuthMiddleware())
	{
		events.GET("/", h.ListEvents)
		events.GET("/:id", h.ListEventByID)
		events.POST("/", h.CreateEvent)
		events.GET("/:id/seats", h.ListAllSeats)
		events.GET("/:id/seatStatus", h.ListSeatsByStatus)
	}
}

func RegisterUserGroup(rg *gin.Engine, h *controllers.UserController) {
	user := rg.Group("/users")
	{
		user.POST("/", h.CreateUser)
	}
}

func RegisterBookingGroup(rg *gin.Engine, h *controllers.BookingController) {
	booking := rg.Group("/booking")
	booking.Use(middlewares.AuthMiddleware())
	{
		booking.POST("/reserve", h.CreateBooking)
		booking.POST("/confirm", h.ConfirmBooking)
		booking.GET("/:id", h.ListBookingByID)
	}
}

func RegisterAuthGroup(rg *gin.Engine, h *controllers.AuthController) {
	auth := rg.Group("/auth")
	{
		auth.POST("/token", h.GenerateToken)
	}
}

func RegisterAdminGroup(rg *gin.Engine, h *controllers.AdminController) {
	admin := rg.Group("/setup-admin")
	{
		admin.POST("/", h.CreateAdmin)
	}
}
