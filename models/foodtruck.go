package models

// JSONFoodTruck is a JSON encodeable version of FoodTruck
type JSONFoodTruck struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Address     string       `json:"address" bson:"address"`
	Location    [2]float64   `json:"location" bson:"location"`
	Owner       string       `json:"owner" bson:"owner"`
	Status      bool         `json:"status" bson:"status"`
	AvgRating   float32      `json:"avgRating" bson:"avgRating"`
	Hours       [7][2]string `json:"hours" bson:"hours"`
	Reviews     []string     `json:"reviews" bson:"reviews"`
	Photos      []string     `json:"photos" bson:"photos"`
	Website     string       `json:"website" bson:"website"`
	PhoneNumber string       `json:"phoneNumber" bson:"phoneNumber"`
	Description string       `json:"description" bson:"description"`
	Tags        []string     `json:"tags" bson:"tags"`
}
