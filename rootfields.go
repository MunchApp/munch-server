package main

import (
	"context"
	"fmt"
	"log"
	"munchserver/models"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Sub-queries underneath root query
func getRootFields() graphql.Fields {
	fields := graphql.Fields{
		"user": &graphql.Field{
			Type: graphql.NewObject(models.GQLUser()),
			Args: graphql.FieldConfigArgument{ //Can define specific user from here
				"email": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email, ok := p.Args["email"].(string)
				var user models.User

				if ok {
					// Connect to MongoDB
					client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
					if err != nil {
						log.Fatal(err)
					}

					db := client.Database("munch")
					userCollection := db.Collection("users")

					// Get specific user from db
					filter := bson.D{{"email", email}}
					err = userCollection.FindOne(context.TODO(), filter).Decode(&user)
					fmt.Println("Found user: %v", user)

					// Disconnect from MongoDB
					err = client.Disconnect(context.TODO())
					if err != nil {
						log.Fatal(err)
						return nil, nil
					}

					return user, nil
				}
				return nil, nil
			},
		},
	}
	return fields
}
