package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Review Model
type Review struct {
	ID       uuid.UUID
	Reviewer *User
	Comment  string
	Rating   float32
	Date     time.Time
}

type JSONReview struct {
	ID       string `json:"_id" bson:"_id"`
	Reviewer string
	Comment  string
	Rating   float32
	Date     string
}

// MarshalJSON encodes a review into JSON
func (review *Review) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSONReview{
		ID:       review.ID.String(),
		Reviewer: review.Reviewer.ID.String(),
		Comment:  review.Comment,
		Rating:   review.Rating,
		Date:     review.Date.String(),
	})
}
