package tests

import (
	"context"
	"munchserver/dbutils"
	"munchserver/models"

	"go.mongodb.org/mongo-driver/mongo"
)

var Db *mongo.Database

func ClearDB() {
	_, _ = Db.Collection("users").DeleteMany(context.TODO(), dbutils.AllQuery())
	_, _ = Db.Collection("foodTrucks").DeleteMany(context.TODO(), dbutils.AllQuery())
	_, _ = Db.Collection("reviews").DeleteMany(context.TODO(), dbutils.AllQuery())
}

func AddFoodTruck(foodTruck models.JSONFoodTruck) {
	_, _ = Db.Collection("foodTrucks").InsertOne(context.TODO(), foodTruck)
}

func AddReview(review models.JSONReview) {
	_, _ = Db.Collection("reviews").InsertOne(context.TODO(), review)
}

func AddUser(user models.JSONUser) {
	_, _ = Db.Collection("users").InsertOne(context.TODO(), user)
}

func GetUser(id string) *models.JSONUser {
	var user models.JSONUser
	err := Db.Collection("users").FindOne(context.TODO(), dbutils.WithIDQuery(id)).Decode(&user)
	if err != nil {
		return nil
	}
	return &user
}

func GetFoodTruck(id string) *models.JSONFoodTruck {
	var foodTruck models.JSONFoodTruck
	err := Db.Collection("foodTrucks").FindOne(context.TODO(), dbutils.WithIDQuery(id)).Decode(&foodTruck)
	if err != nil {
		return nil
	}
	return &foodTruck
}

func GetReview(id string) *models.JSONReview {
	var review models.JSONReview
	err := Db.Collection("reviews").FindOne(context.TODO(), dbutils.WithIDQuery(id)).Decode(&review)
	if err != nil {
		return nil
	}
	return &review
}
