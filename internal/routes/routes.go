package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/controllers"
)

func RegisterEventsGroup(rg *gin.Engine, h *controllers.EventController) {
	events := rg.Group("/events")
	{
		events.GET("/", h.ListEvents)
		events.GET("/:id", h.ListEventByID)
		events.POST("/", h.CreateEvent)
		events.GET("/:id/seats", h.ListAllSeatsAvailables)

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
	{
		booking.POST("/", h.CreateBooking)
	}
}
