package models

import (
	"time"
)

type JSONReview struct {
	ID           string    `json:"id" bson:"_id"`
	Reviewer     string    `json:"reviewer" bson:"reviewer"`
	ReviewerName string    `json:"reviewerName" bson:"reviewerName"`
	FoodTruck    string    `json:"foodTruck" bson:"foodTruck"`
	Comment      string    `json:"comment" bson:"comment"`
	Rating       float32   `json:"rating" bson:"rating"`
	Date         time.Time `json:"date" bson:"date"`
	Origin       string    `json:"origin" bson:"origin"`
}
