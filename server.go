package main

import (
	"context"
	"fmt"
	"log"
	"munchserver/routes"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Setup http router
	router := mux.NewRouter()
	router.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	router.HandleFunc("/foodtrucks", routes.GetFoodTrucksHandler).Methods("GET")
	router.HandleFunc("/reviews", routes.GetReviewsHandler).Methods("GET")
	router.HandleFunc("/contributors", routes.GetContributorsHandler).Methods("GET")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))

	if err != nil {
		log.Fatal(err)
		return
	}

	db := client.Database("munch")

	routes.Db = db
	routes.Router = router

	fmt.Println("Connected to MongoDB!")
	log.Fatal(http.ListenAndServe(":80", router))
}
