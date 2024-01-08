package controllers

import (
	"app/busstationgo/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ticketsCollection *mongo.Collection

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connectionString := os.Getenv("MONGODB_CONNECTION_STRING")
	clientOptions := options.Client().ApplyURI(connectionString)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mongodb connection success")

	dbName := os.Getenv("DB_NAME")
	colName := "tickets"

	ticketsCollection = client.Database(dbName).Collection(colName)

	fmt.Println("Collection instance is ready")
}

func BuyTicket(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var ticket models.BuyTicketVM

	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetHeader("x-user-id")

	userIdInt, err := strconv.Atoi(userId)

	newTicket := models.Ticket{
		RouteId:   ticket.RouteId,
		SeatNo:    ticket.SeatNo,
		UserId:    userIdInt,
		Confirmed: true,
		Date:      primitive.NewDateTimeFromTime(time.Now()),
	}

	fmt.Println(newTicket.Date)

	count, err := ticketsCollection.CountDocuments(ctx, bson.M{"userId": newTicket.UserId, "routeId": ticket.RouteId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking ticket"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket already bought"})
		return
	}

	insertedTicket, insertErr := ticketsCollection.InsertOne(ctx, newTicket)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while buying ticket"})
		return
	}

	newTicket.ID = insertedTicket.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, newTicket)
}

func BookTicket(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	var ticket models.BuyTicketVM

	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetHeader("x-user-id")

	userIdInt, err := strconv.Atoi(userId)

	newTicket := models.Ticket{
		RouteId:   ticket.RouteId,
		SeatNo:    ticket.SeatNo,
		UserId:    userIdInt,
		Confirmed: false,
		Date:      primitive.NewDateTimeFromTime(time.Now()),
	}

	count, err := ticketsCollection.CountDocuments(ctx, bson.M{"userId": newTicket.UserId, "routeId": ticket.RouteId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking ticket"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket already exists"})
		return
	}

	insertedTicket, insertErr := ticketsCollection.InsertOne(ctx, newTicket)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while persisting ticket"})
		return
	}

	newTicket.ID = insertedTicket.InsertedID.(primitive.ObjectID)

	c.JSON(http.StatusOK, newTicket)
}

func ConfirmTicket(c *gin.Context) {
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	tid := c.Param("id")

	ticketId, err := primitive.ObjectIDFromHex(tid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket id"})
		return
	}

	var ticketToConfirm models.Ticket

	userId := c.GetHeader("x-user-id")

	userIdInt, err := strconv.Atoi(userId)

	err = ticketsCollection.FindOne(ctx, bson.M{"_id": ticketId, "userId": userIdInt}).Decode(&ticketToConfirm)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket not found"})
		return
	}

	if ticketToConfirm.Confirmed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket already confirmed"})
		return
	}

	ticketToConfirm.Confirmed = true

	_, updateErr := ticketsCollection.UpdateOne(ctx, bson.M{"_id": ticketId}, bson.M{"$set": ticketToConfirm})
	if updateErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while confirming ticket"})
		return
	}

	c.JSON(http.StatusOK, ticketToConfirm)
}

func CheckExpiredReservations() {
	// find all tickets that are not confirmed and are older than 7 days
	var ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	cursor, err := ticketsCollection.Find(ctx, bson.M{"confirmed": false, "date": bson.M{"$lt": primitive.NewDateTimeFromTime(time.Now().AddDate(0, 0, -7))}})
	if err != nil {
		log.Fatal(err)
	}

	var tickets []models.Ticket
	if err = cursor.All(ctx, &tickets); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Found", len(tickets), "expired reservations")

	for _, ticket := range tickets {
		_, err := ticketsCollection.DeleteOne(ctx, bson.M{"_id": ticket.ID})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Expired reservations deleted")
}
