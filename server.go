package main

import (
	"context"
	"fmt"
	"log"
	"munchserver/middleware"
	"munchserver/models"
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
	router.HandleFunc("/reviews", routes.GetReviewsHandler).Methods("GET")
	router.HandleFunc("/contributors", routes.GetContributorsHandler).Methods("GET")

	// Auth required routes
	router.Use(middleware.AuthenticateUser)
	router.HandleFunc("/foodtrucks", routes.PostFoodTrucksHandler).Methods("POST")
	router.HandleFunc("/reviews", routes.PostReviewsHandler).Methods("POST")

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
	models.Db = db

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
	log.Fatal(http.ListenAndServe(":"+secrets.GetPort(), router))
}
