package main

import (
	"app/busstationgo/controllers"
	"app/busstationgo/routers"
	"fmt"
	"log"
	"net/http"

	"github.com/robfig/cron"
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

	r := gin.Default()

	ticketRouter := routers.TicketsRouter()

	ticketGroup := r.Group("/tickets")
	ticketGroup.Any("/*path", gin.WrapH(ticketRouter))

	healthcheck.New(r, config.DefaultConfig(), []checks.Check{})

	c := cron.New()

	c.AddFunc("@every 5m", func() {
		controllers.CheckExpiredReservations()
	})

	c.Start()

	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
