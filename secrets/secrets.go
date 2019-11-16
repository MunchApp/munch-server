package secrets

import (
	"log"
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

func GetAWSAccessKey() string {
	accessKey, exists := os.LookupEnv("AWS_ACCESS_KEY")
	if !exists {
		log.Println("AWS Access key not found, s3 upload will not work")
	}
	return accessKey
}

func GetAWSSecretAccessKey() string {
	secretAccessKey, exists := os.LookupEnv("AWS_SECRET_ACCESS_KEY")
	if !exists {
		log.Println("AWS Secret access key not found, s3 upload will not work")
	}
	return secretAccessKey
}
