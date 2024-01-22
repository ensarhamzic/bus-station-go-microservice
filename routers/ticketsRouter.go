package routers

import (
	"app/busstationgo/controllers"

	"github.com/gin-gonic/gin"
)

func TicketsRouter() *gin.Engine {
	router := gin.Default()

	ticketsGroup := router.Group(TicketsBaseRoute)
	{
		ticketsGroup.POST("/buy", controllers.BuyTicket)
		ticketsGroup.POST("/book", controllers.BookTicket)
		ticketsGroup.POST("/confirm/:id", controllers.ConfirmTicket)
		ticketsGroup.GET("/check/:routeId/:seatNo", controllers.CheckTicketAvailability)
	}

	return router
}
