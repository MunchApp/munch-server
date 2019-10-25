package models

import (
	"encoding/json"

	"github.com/google/uuid"
)

// FoodTruck Model
type FoodTruck struct {
	ID          uuid.UUID
	Name        string
	Address     string
	Location    [2]float64
	Owner       *User
	Status      bool
	AvgRating   float32
	Hours       [7][2]string
	Reviews     []*Review
	Photos      []string
	Website     string
	PhoneNumber string
	Description string
	Tags        []string
}

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

func NewJSONFoodTruck(foodTruck FoodTruck) JSONFoodTruck {
	reviews := make([]string, len(foodTruck.Reviews))
	for i, review := range foodTruck.Reviews {
		reviews[i] = review.ID.String()
	}
	return JSONFoodTruck{
		ID:          foodTruck.ID.String(),
		Name:        foodTruck.Name,
		Address:     foodTruck.Address,
		Location:    foodTruck.Location,
		Owner:       foodTruck.Owner.ID.String(),
		Status:      foodTruck.Status,
		AvgRating:   foodTruck.AvgRating,
		Hours:       foodTruck.Hours,
		Reviews:     reviews,
		Photos:      foodTruck.Photos,
		Website:     foodTruck.Website,
		PhoneNumber: foodTruck.PhoneNumber,
		Description: foodTruck.Description,
		Tags:        foodTruck.Tags,
	}
}

// MarshalJSON encodes a food truck into JSON
func (foodTruck FoodTruck) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONFoodTruck(foodTruck))
}
