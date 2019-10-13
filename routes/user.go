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
	usersCollection := Db.Collection("users")
	cur, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}
	var users []models.JSONUser
	cur.All(context.TODO(), &users)

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
