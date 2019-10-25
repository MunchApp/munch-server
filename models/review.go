package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Review Model
type Review struct {
	ID           uuid.UUID
	Reviewer     *User
	ReviewerName string
	FoodTruck    *FoodTruck
	Comment      string
	Rating       float32
	Date         time.Time
	Origin       string
}

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

func NewJSONReview(review Review) JSONReview {
	return JSONReview{
		ID:           review.ID.String(),
		Reviewer:     review.Reviewer.ID.String(),
		ReviewerName: review.ReviewerName,
		FoodTruck:    review.FoodTruck.ID.String(),
		Comment:      review.Comment,
		Rating:       review.Rating,
		Date:         review.Date,
		Origin:       review.Origin,
	}
}

// MarshalJSON encodes a review into JSON
func (review Review) MarshalJSON() ([]byte, error) {
	return json.Marshal(NewJSONReview(review))
}
