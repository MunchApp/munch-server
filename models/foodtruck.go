package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FoodTruck Model
type FoodTruck struct {
	ID          uuid.UUID
	Name        string
	Address     string
	Owner       *User
	AvgRating   float32
	Hours       [2]time.Time
	Reviews     []*Review
	Photos      []string
	Website     string
	PhoneNumber string
}

// JSONFoodTruck is a JSON encodable version of FoodTruck
type JSONFoodTruck struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Address     string       `json:"address" bson:"address"`
	Owner       string       `json:"owner" bson:"owner"`
	AvgRating   float32      `json:"avgRating" bson:"avgRating"`
	Hours       [2]time.Time `json:"hours" bson:"hours"`
	Reviews     []string     `json:"reviews" bson:"reviews"`
	Photos      []string     `json:"photos" bson:"photos"`
	Website     string       `json:"website" bson:"website"`
	PhoneNumber string       `json:"phoneNumber" bson:"phoneNumber"`
}

func NewJSONFoodTruck(foodTruck FoodTruck) JSONFoodTruck {
	reviews := make([]string, len(foodTruck.Reviews))
	for i, review := range foodTruck.Reviews {
		reviews[i] = review.ID.String()
	}
	return JSONFoodTruck{
		ID:        foodTruck.ID.String(),
		Name:      foodTruck.Name,
		Address:   foodTruck.Address,
		Owner:     foodTruck.Owner.ID.String(),
		AvgRating: foodTruck.AvgRating,
		Hours:     foodTruck.Hours,
		Reviews:   reviews,
	}
}

// MarshalJSON encodes a food truck into JSON
func (foodTruck FoodTruck) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONFoodTruck(foodTruck))
}
