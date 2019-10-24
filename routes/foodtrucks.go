package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"munchserver/models"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

type addFoodTruckRequest struct {
	Name        *string       `json:"name"`
	Address     *string       `json:"address"`
	Hours       *[2]time.Time `json:"hours"`
	Photos      *[]string     `json:"photos"`
	Website     string        `json:"website"`
	PhoneNumber string        `json:"phoneNumber"`
}

func PostFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {
	foodTruckDecoder := json.NewDecoder(r.Body)
	foodTruckDecoder.DisallowUnknownFields()

	var newFoodTruck addFoodTruckRequest
	err := foodTruckDecoder.Decode(&newFoodTruck)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

	addFoodTruck := models.JSONFoodTruck{
		ID:          uuid.String(),
		Name:        *newFoodTruck.Name,
		Address:     *newFoodTruck.Address,
		Hours:       *newFoodTruck.Hours,
		Photos:      *newFoodTruck.Photos,
		Website:     newFoodTruck.Website,
		PhoneNumber: newFoodTruck.PhoneNumber,
	}
	_, err = Db.Collection("foodTrucks").InsertOne(context.TODO(), addFoodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
