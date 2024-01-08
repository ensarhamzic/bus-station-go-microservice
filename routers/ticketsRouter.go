package routers

import (
	"app/busstationgo/controllers"

	"github.com/gin-gonic/gin"
)

func TicketsRouter() *gin.Engine {
	router := gin.Default()

	ticketsGroup := router.Group("/tickets")
	{
		ticketsGroup.POST("/buy", controllers.BuyTicket)
		ticketsGroup.POST("/book", controllers.BookTicket)
		ticketsGroup.POST("/confirm/:id", controllers.ConfirmTicket)
	}

	return router
}
