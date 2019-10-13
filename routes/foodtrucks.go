package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"munchserver/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {
	// Get all users from the database into a cursor
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
