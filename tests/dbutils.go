package tests

import (
	"context"
	"munchserver/queries"

	"go.mongodb.org/mongo-driver/mongo"
)

var db *mongo.Database

func clearDB() {
	db.Collection("users").DeleteMany(context.TODO(), queries.All())
	db.Collection("foodTrucks").DeleteMany(context.TODO(), queries.All())
	db.Collection("reviews").DeleteMany(context.TODO(), queries.All())
}
