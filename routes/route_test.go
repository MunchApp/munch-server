package routes

import (
	"context"
	"munchserver/secrets"
	"munchserver/tests"
	"os"
	"testing"

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
	code := m.Run()

	os.Exit(code)
}
