package tests

import (
	"context"
	"munchserver/models"
	"munchserver/queries"

	"go.mongodb.org/mongo-driver/mongo"
)

var Db *mongo.Database

func ClearDB() {
	Db.Collection("users").DeleteMany(context.TODO(), queries.All())
	Db.Collection("foodTrucks").DeleteMany(context.TODO(), queries.All())
	Db.Collection("reviews").DeleteMany(context.TODO(), queries.All())
}

func AddFoodTruck(foodTruck models.JSONFoodTruck) {
	Db.Collection("foodTrucks").InsertOne(context.TODO(), foodTruck)
}

func AddReview(review models.JSONReview) {
	Db.Collection("reviews").InsertOne(context.TODO(), review)
}

func AddUser(user models.JSONUser) {
	Db.Collection("users").InsertOne(context.TODO(), user)
}
