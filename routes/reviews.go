package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"munchserver/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all users from the database into a cursor
	reviewsCollection := Db.Collection("reviews")
	cur, err := reviewsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}

	// Get users from cursor, convert to empty slice if no users in DB
	var reviews []models.JSONReview
	cur.All(context.TODO(), &reviews)
	if reviews == nil {
		reviews = make([]models.JSONReview, 0)
	}

	// Convert users to json
	js, err := json.Marshal(reviews)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in decoding mongo document: %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
