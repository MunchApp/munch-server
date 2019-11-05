package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"munchserver/middleware"
	"munchserver/models"
	"munchserver/queries"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type newReviewRequest struct {
	ReviewerName string    `json:"reviewerName"`
	FoodTruck    *string   `json:"foodTruck"`
	Comment      string    `json:"comment"`
	Rating       *float32  `json:"rating"`
	Date         time.Time `json:"date"`
	Origin       string    `json:"origin"`
}

func PostReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user, or if the user agent is from the scraper
	if !userLoggedIn && r.Header.Get("User-Agent") != "MunchCritic/1.0" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	reviewDecoder := json.NewDecoder(r.Body)
	reviewDecoder.DisallowUnknownFields()

	// Decode request
	var newReview newReviewRequest
	err := reviewDecoder.Decode(&newReview)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure required fields set
	if (!userLoggedIn && newReview.ReviewerName == "") ||
		newReview.FoodTruck == nil ||
		newReview.Rating == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Lookup food truck

	// Generate uuid for food truck
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	date := newReview.Date
	if date.IsZero() {
		date = time.Now()
	}
	origin := newReview.Origin
	if origin == "" {
		origin = "munchapp"
	}
	reviewer := ""
	if userLoggedIn {
		reviewer = user
	}

	addedReview := models.JSONReview{
		ID:           uuid.String(),
		Reviewer:     reviewer,
		ReviewerName: newReview.ReviewerName,
		FoodTruck:    *newReview.FoodTruck,
		Comment:      newReview.Comment,
		Rating:       *newReview.Rating,
		Date:         date,
		Origin:       origin,
	}

	// Add review to database
	_, err = Db.Collection("reviews").InsertOne(r.Context(), addedReview)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Attach review to user
	if userLoggedIn {
		_, err = Db.Collection("users").UpdateOne(r.Context(), queries.WithID(user), queries.PushReview(uuid.String()))
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Attach review to food truck
	_, err = Db.Collection("foodTrucks").UpdateOne(r.Context(), queries.WithID(*newReview.FoodTruck), queries.PushReview(uuid.String()))
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

func GetReviewsOfFoodTruckHandler(w http.ResponseWriter, r *http.Request) {
	// Get food truck id from route params
	params := mux.Vars(r)
	foodTruckID, foodTruckIDExists := params["foodTruckID"]

	if !foodTruckIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get all reviews with foodtruck from the database into a cursor
	reviewsCollection := Db.Collection("reviews")
	cur, err := reviewsCollection.Find(context.TODO(), queries.WithFoodTruck(foodTruckID))
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}

	// Get reviews from cursor, convert to empty slice if no reviews in DB
	var reviews []models.JSONReview
	cur.All(context.TODO(), &reviews)
	if reviews == nil {
		reviews = make([]models.JSONReview, 0)
	}

	// Convert reviews to json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func GetReviewsHandler(w http.ResponseWriter, r *http.Request) {
	// Get all reviews from the database into a cursor
	reviewsCollection := Db.Collection("reviews")
	cur, err := reviewsCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}

	// Get reviews from cursor, convert to empty slice if no reviews in DB
	var reviews []models.JSONReview
	cur.All(context.TODO(), &reviews)
	if reviews == nil {
		reviews = make([]models.JSONReview, 0)
	}

	// Convert reviews to json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}
