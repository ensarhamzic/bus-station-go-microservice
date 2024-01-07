package routers

import (
	"app/busstationgo/controllers"

	"github.com/gin-gonic/gin"
)

func TicketsRouter() *gin.Engine {
	router := gin.Default()

	ticketsGroup := router.Group("/tickets")
	{
		ticketsGroup.POST("/", controllers.BuyTicket)
	}

	return router
}
