package models

import (
	"encoding/json"
	"net/http"
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
	ID       string  `json:"_id"`
	Reviewer string  `json:"reviewer"`
	Comment  string  `json:"comment"`
	Rating   float32 `json:"rating"`
	Date     string  `json:"date"`
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

func ReviewHandler(w http.ResponseWriter, r *http.Request) {

}
