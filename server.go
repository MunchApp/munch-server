package main

import (
	"context"
	"fmt"
	"log"
<<<<<<< HEAD
=======
	"munchserver/gqlfields"
	_ "munchserver/gqlfields"
	_ "munchserver/models"
>>>>>>> working root query, finished mutation query code but not working
	"munchserver/routes"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Setup http router
	router := mux.NewRouter()
	router.HandleFunc("/register", routes.PostRegisterHandler).Methods("POST")
	router.HandleFunc("/login", routes.PostLoginHandler).Methods("POST")
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

	// Inject db to routes
	routes.Db = db
	routes.Router = router

	// Find port from env var or default to 80
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	// Setup db indexes
	userIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true).SetBackground(true),
	}
	_, err = db.Collection("users").Indexes().CreateOne(context.TODO(), userIndex)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	log.Fatal(http.ListenAndServe(":"+port, router))
}
