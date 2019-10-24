package main

import (
	"context"
	"fmt"
	"log"
	"munchserver/routes"
	"net/http"
	"os"

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
	mongoURI, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		mongoURI = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))

	if err != nil {
		log.Fatal(err)
		return
	}

	dbName, exists := os.LookupEnv("MONGODB_DBNAME")
	if !exists {
		dbName = "munch"
	}

	db := client.Database(dbName)

	routes.Db = db
	routes.Router = router

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	fmt.Println("Connected to MongoDB!")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
