package models

import (
	"encoding/json"
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
	ID              string `json:"_id" bson:"_id"`
	Name            string
	Favorites       []string
	Reviews         []string
	Email           string
	OwnedFoodTrucks []string
}

// MarshalJSON encodes a user into JSON
func (user *User) MarshalJSON() ([]byte, error) {
	favorites := make([]string, len(user.Favorites))
	reviews := make([]string, len(user.Reviews))
	ownedFoodTrucks := make([]string, len(user.OwnedFoodTrucks))
	for i, favorite := range user.Favorites {
		favorites[i] = favorite.ID.String()
	}
	for i, review := range user.Reviews {
		reviews[i] = review.ID.String()
	}
	for i, ownedFoodTruck := range user.OwnedFoodTrucks {
		ownedFoodTrucks[i] = ownedFoodTruck.ID.String()
	}
	return json.Marshal(JSONUser{
		ID:              user.ID.String(),
		Name:            user.Name,
		Favorites:       favorites,
		Reviews:         reviews,
		Email:           user.Email,
		OwnedFoodTrucks: ownedFoodTrucks,
	})
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

}
