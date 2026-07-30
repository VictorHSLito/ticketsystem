package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/victorhslito/ticketsystem/internal/controllers"
)

func RegisterEventsGroup(rg *gin.Engine) {
	events := rg.Group("/events")
	{
		events.GET("/", controllers.ListEvents)
		events.GET("/:id", controllers.ListEventByID)
		events.POST("/", controllers.CreateEvent)
		events.GET("/:id/seats", controllers.ListAllSeatsAvailables)

	}
}

func RegisterUserGroup(rg *gin.Engine) {
	user := rg.Group("/users")
	{
		user.POST("/", controllers.CreateUser)
	}
}

func RegisterBookingGroup(rg *gin.Engine) {
	booking := rg.Group("/booking")
	{
		booking.POST("/", controllers.CreateBooking)
	}
}
