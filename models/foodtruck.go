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
	ID        string `json:"_id" bson:"_id"`
	Name      string
	Address   string
	Owner     string
	AvgRating float32
	Hours     []string
	Reviews   []string
}

// MarshalJSON encodes a food truck into JSON
func (foodTruck *FoodTruck) MarshalJSON() ([]byte, error) {
	hours := []string{foodTruck.Hours[0].String(), foodTruck.Hours[1].String()}
	reviews := make([]string, len(foodTruck.Reviews))
	for i, review := range foodTruck.Reviews {
		reviews[i] = review.ID.String()
	}
	return json.Marshal(JSONFoodTruck{
		ID:        foodTruck.ID.String(),
		Name:      foodTruck.Name,
		Address:   foodTruck.Address,
		Owner:     foodTruck.Owner.ID.String(),
		AvgRating: foodTruck.AvgRating,
		Hours:     hours,
		Reviews:   reviews,
	})
}
