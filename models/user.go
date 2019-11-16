package models

import (
	"time"
)

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
