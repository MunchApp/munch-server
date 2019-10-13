package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"munchserver/models"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get all users from the database into a cursor
	usersCollection := Db.Collection("users")
	cur, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}

	// Get users from cursor, convert to empty slice if no users in DB
	var users []models.JSONUser
	cur.All(context.TODO(), &users)
	if users == nil {
		users = make([]models.JSONUser, 0)
	}

	// Convert users to json
	js, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in decoding mongo document: %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
