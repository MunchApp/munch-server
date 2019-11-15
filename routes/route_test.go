package routes

import (
	"context"
	"log"
	"munchserver/secrets"
	"munchserver/tests"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(secrets.GetMongoURI()))
	if err != nil {
		panic(err)
	}

	Db = client.Database(secrets.GetTestMongoDBName())
	// Inject db to tests
	tests.Db = Db

	tests.ClearDB()

	// Setup db indexes
	userIndex := mongo.IndexModel{
		Keys:    bson.M{"email": 1},
		Options: options.Index().SetUnique(true).SetBackground(true),
	}
	_, err = Db.Collection("users").Indexes().CreateOne(context.TODO(), userIndex)
	if err != nil {
		log.Fatal(err)
	}

	code := m.Run()

	tests.ClearDB()

	os.Exit(code)
}
