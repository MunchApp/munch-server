package models

import (
	"net/http"

	"github.com/google/uuid"
)

// User Model
type User struct {
	ID              uuid.UUID
	Name            string
	Favorites       []*FoodTruck
	Reviews         []*Review
	Email           string
	OwnedFoodTrucks []*FoodTruck
}

type JSONUser struct {
	ID              string   `json:"_id"`
	Name            string   `json:"name"`
	Favorites       []string `json:"favorites"`
	Reviews         []string `json:"reviews"`
	Email           string   `json:"email"`
	OwnedFoodTrucks []string `json:"ownedFoodTrucks"`
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

}
