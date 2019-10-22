package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// FoodTruck Model
type FoodTruck struct {
	ID        uuid.UUID
	Name      string
	Address   string
	Owner     *User
	AvgRating float32
	Hours     [2]time.Time
	Reviews   []*Review
}

// JSONFoodTruck is a JSON encodable version of FoodTruck
type JSONFoodTruck struct {
	ID        string   `json:"id" bson:"_id"`
	Name      string   `json:"name" bson:"name"`
	Address   string   `json:"address" bson:"address"`
	Owner     string   `json:"owner" bson:"owner"`
	AvgRating float32  `json:"avgRating" bson:"avgRating"`
	Hours     []string `json:"hours" bson:"hours"`
	Reviews   []string `json:"reviews" bson:"reviews"`
}

func NewJSONFoodTruck(foodTruck FoodTruck) JSONFoodTruck {
	hours := []string{foodTruck.Hours[0].String(), foodTruck.Hours[1].String()}
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
		Hours:     hours,
		Reviews:   reviews,
	}
}

// MarshalJSON encodes a food truck into JSON
func (foodTruck FoodTruck) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONFoodTruck(foodTruck))
}
