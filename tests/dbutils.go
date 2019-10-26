package tests

import (
	"context"
	"munchserver/queries"

	"go.mongodb.org/mongo-driver/mongo"
)

var Db *mongo.Database

func ClearDB() {
	Db.Collection("users").DeleteMany(context.TODO(), queries.All())
	Db.Collection("foodTrucks").DeleteMany(context.TODO(), queries.All())
	Db.Collection("reviews").DeleteMany(context.TODO(), queries.All())
}
