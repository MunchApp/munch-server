package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// User Model
type User struct {
	ID              uuid.UUID
	PasswordHash    []byte
	NameFirst       string
	NameLast        string
	Email           string
	PhoneNumber     string
	City            string
	State           string
	DateOfBirth     time.Time
	OwnedFoodTrucks []*FoodTruck
	Favorites       []*FoodTruck
	Reviews         []*Review
}

type JSONUser struct {
	ID              string    `json:"id" bson:"_id"`
	PasswordHash    []byte    `json:"passwordHash" bson:"passwordHash"`
	NameFirst       string    `json:"firstName" bson:"firstName"`
	NameLast        string    `json:"lastName" bson:"lastName"`
	Email           string    `json:"email" bson:"email"`
	PhoneNumber     string    `json:"phoneNumber" bson:"phoneNumber"`
	City            string    `json:"city" bson:"city"`
	State           string    `json:"state" bson:"state"`
	DateOfBirth     time.Time `json:"dateOfBirth" bson:"dateOfBirth"`
	Favorites       []string  `json:"favorites" bson:"favorites"`
	Reviews         []string  `json:"reviews" bson:"reviews"`
	OwnedFoodTrucks []string  `json:"ownedFoodTrucks" bson:"ownedFoodTrucks"`
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
		PasswordHash:    user.PasswordHash,
		NameFirst:       user.NameFirst,
		NameLast:        user.NameLast,
		Email:           user.Email,
		PhoneNumber:     user.PhoneNumber,
		City:            user.City,
		State:           user.State,
		DateOfBirth:     user.DateOfBirth,
		Favorites:       favorites,
		Reviews:         reviews,
		OwnedFoodTrucks: ownedFoodTrucks,
	})
}
