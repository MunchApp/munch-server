package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"munchserver/middleware"
	"munchserver/models"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type addFoodTruckRequest struct {
	Name        *string       `json:"name"`
	Address     *string       `json:"address"`
	Location    [2]float64    `json:"location"`
	Hours       *[7][2]string `json:"hours"`
	Photos      *[]string     `json:"photos"`
	Website     string        `json:"website"`
	PhoneNumber string        `json:"phoneNumber"`
	Description string        `json:"description"`
	Tags        []string      `json:"tags"`
}

func PostFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user := r.Context().Value(middleware.UserKey)

	// Check for a user, or if the user agent is from the scraper
	if user == nil && r.Header.Get("User-Agent") != "MunchCritic/1.0" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	foodTruckDecoder := json.NewDecoder(r.Body)
	foodTruckDecoder.DisallowUnknownFields()

	var newFoodTruck addFoodTruckRequest
	err := foodTruckDecoder.Decode(&newFoodTruck)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure required fields are set
	if newFoodTruck.Name == nil ||
		newFoodTruck.Address == nil ||
		newFoodTruck.Hours == nil ||
		newFoodTruck.Photos == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set tags to an empty array if they don't exist
	tags := newFoodTruck.Tags
	if tags == nil {
		tags = []string{}
	}

	addedFoodTruck := models.JSONFoodTruck{
		ID:          uuid.String(),
		Name:        *newFoodTruck.Name,
		Address:     *newFoodTruck.Address,
		Location:    newFoodTruck.Location,
		Owner:       user.(string),
		Hours:       *newFoodTruck.Hours,
		Reviews:     []string{},
		Photos:      *newFoodTruck.Photos,
		Website:     newFoodTruck.Website,
		PhoneNumber: newFoodTruck.PhoneNumber,
		Description: newFoodTruck.Description,
		Tags:        tags,
	}

	// Add food truck to database
	_, err = Db.Collection("foodTrucks").InsertOne(context.TODO(), addedFoodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

func GetFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {
	// Get all foodtrucks from the database into a cursor
	foodTrucksCollection := Db.Collection("foodTrucks")
	cur, err := foodTrucksCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}

	// Get users from cursor, convert to empty slice if no users in DB
	var foodTrucks []models.JSONFoodTruck
	cur.All(context.TODO(), &foodTrucks)
	if foodTrucks == nil {
		foodTrucks = make([]models.JSONFoodTruck, 0)
	}

	// Convert users to json
	js, err := json.Marshal(foodTrucks)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in decoding mongo document: %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
