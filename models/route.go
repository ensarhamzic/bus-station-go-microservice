package models

type Route struct {
	ID            int     `json:"_id,omitempty" bson:"_id,omitempty"`
	BusId         int     `json:"busId,omitempty" bson:"busId,omitempty"`
	DriverId      int     `json:"driverId,omitempty" bson:"driverId,omitempty"`
	FromLocation  int     `json:"fromLocation,omitempty" bson:"fromLocation,omitempty"`
	ToLocation    int     `json:"toLocation,omitempty" bson:"toLocation,omitempty"`
	DepartureTime int64   `json:"departureTime,omitempty" bson:"departureTime,omitempty"`
	ArrivalTime   int64   `json:"arrivalTime,omitempty" bson:"arrivalTime,omitempty"`
	Price         float64 `json:"price,omitempty" bson:"price,omitempty"`
}
