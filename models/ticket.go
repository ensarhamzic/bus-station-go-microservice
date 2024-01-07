package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId  int                `json:"userId,omitempty" bson:"userId,omitempty"`
	RouteId int                `json:"routeId,omitempty" bson:"routeId,omitempty"`
	SeatNo  int                `json:"seatNo,omitempty" bson:"seatNo,omitempty"`
}
