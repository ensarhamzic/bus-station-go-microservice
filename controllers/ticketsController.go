package controllers

import (
	"app/busstationgo/models"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

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

	var ticket models.Ticket

	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	count, err := ticketsCollection.CountDocuments(ctx, bson.M{"userId": ticket.UserId, "routeId": ticket.RouteId})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking ticket"})
		return
	}

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ticket already bought"})
		return
	}

	insertedTicket, insertErr := ticketsCollection.InsertOne(ctx, ticket)

	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while buying ticket"})
		return
	}

	ticket.ID = insertedTicket.InsertedID.(primitive.ObjectID)
	c.JSON(http.StatusOK, ticket)
}
