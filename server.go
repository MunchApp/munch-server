package munchserver

import (
	"fmt"
	"log"
	"munchserver/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	db     *mongo.Database
	router *mux.Router
)

func main() {
	// Setup http router
	router = mux.NewRouter()
	router.HandleFunc("/users", models.UserHandler)
	router.HandleFunc("/foodtrucks", models.FoodTruckHandler)
	router.HandleFunc("/reviews", models.ReviewHandler)

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
	}

	db = client.Database("munch")

	fmt.Println("Connected to MongoDB!")
}
