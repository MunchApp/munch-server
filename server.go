package main

import (
	"context"
	"fmt"
	"log"
	"munchserver/middleware"
	"munchserver/routes"
	"munchserver/secrets"
	"net/http"

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
	router.HandleFunc("/foodtrucks/{foodTruckID}", routes.GetFoodTruckHandler).Methods("GET")
	router.HandleFunc("/reviews", routes.GetReviewsHandler).Methods("GET")
	router.HandleFunc("/reviews/{reviewID}", routes.GetReviewHandler).Methods("GET")
	router.HandleFunc("/reviews/foodtruck/{foodTruckID}", routes.GetReviewsOfFoodTruckHandler).Methods("GET")
	router.HandleFunc("/contributors", routes.GetContributorsHandler).Methods("GET")
	router.HandleFunc("/users/{userID}", routes.GetUserHandler).Methods("GET")

	// Auth required routes
	router.Use(middleware.AuthenticateUser)
	router.HandleFunc("/profile", routes.GetProfileHandler).Methods("GET")
	router.HandleFunc("/foodtrucks", routes.PostFoodTrucksHandler).Methods("POST")
	router.HandleFunc("/foodtrucks/claim/{foodTruckID}", routes.PutClaimFoodTruckHandler).Methods("PUT")
	router.HandleFunc("/reviews", routes.PostReviewsHandler).Methods("POST")
	router.HandleFunc("/users/favorite/{foodTruckID}", routes.PutFavoriteHandler).Methods("PUT")
	router.HandleFunc("/profile", routes.PutUpdateProfileHandler).Methods("PUT")
	router.HandleFunc("/foodtrucks/{foodTruckID}", routes.PutFoodTrucksHandler).Methods("PUT")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(secrets.GetMongoURI()))

	if err != nil {
		log.Fatal(err)
		return
	}

	db := client.Database(secrets.GetMongoDBName())

	// Inject db to routes
	routes.Db = db
	routes.Router = router

	// Setup db indexes
	userIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true).SetBackground(true),
	}
	_, err = db.Collection("users").Indexes().CreateOne(context.TODO(), userIndex)
	if err != nil {
		log.Fatal(err)
	}
	locationIndex := mongo.IndexModel{
		Keys: bson.M{"location": "2dsphere"},
	}
	_, err = db.Collection("foodTrucks").Indexes().CreateOne(context.TODO(), locationIndex)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")
	log.Fatal(http.ListenAndServe(":"+secrets.GetPort(), router))
}
