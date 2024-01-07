package main

import (
	"app/busstationgo/routers"
	"fmt"
	"log"
	"net/http"

	healthcheck "github.com/tavsec/gin-healthcheck"
	"github.com/tavsec/gin-healthcheck/checks"
	"github.com/tavsec/gin-healthcheck/config"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("Server is getting started...")

	router := gin.Default()

	ticketRouter := routers.TicketsRouter()

	ticketGroup := router.Group("/tickets")
	ticketGroup.Any("/", gin.WrapH(ticketRouter))

	healthcheck.New(router, config.DefaultConfig(), []checks.Check{})

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
