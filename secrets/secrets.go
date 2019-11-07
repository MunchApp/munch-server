package secrets

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

// TODO: Move the secret to an env var
func GetJWTSecret(token *jwt.Token) (interface{}, error) {
	return []byte("MunchIsReallyCool"), nil
}

func GetMongoURI() string {
	mongoURI, exists := os.LookupEnv("MONGODB_URI")
	if !exists {
		mongoURI = "mongodb://localhost:27017"
	}
	return mongoURI
}

func GetMongoDBName() string {
	dbName, exists := os.LookupEnv("MONGODB_DBNAME")
	if !exists {
		dbName = "munch"
	}
	return dbName
}

func GetTestMongoDBName() string {
	dbName, exists := os.LookupEnv("MONGODB_TESTDBNAME")
	if !exists {
		dbName = "munch_test"
	}
	return dbName
}

func GetPort() string {
	// Find port from env var or default to 80
	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}
	return port
}
