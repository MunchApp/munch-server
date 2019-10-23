package gqlfields

import (
	"context"
	"fmt"
	"log"
	"munchserver/models"
	"munchserver/routes"

	"github.com/graphql-go/graphql"
	"go.mongodb.org/mongo-driver/bson"
)

// Sub-queries underneath root query
func GetRootQueryFields() graphql.Fields {
	fields := graphql.Fields{
		// Get user info
		"user": &graphql.Field{
			Type: UserType,
			Args: graphql.FieldConfigArgument{ //Can define specific user from here
				"email": &graphql.ArgumentConfig{Type: graphql.String},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email, ok := p.Args["email"].(string)
				var user models.User

				if ok {

					userCollection := routes.Db.Collection("users")

					// Get specific user from db
					filter := bson.D{{"email", email}}
					err := userCollection.FindOne(context.TODO(), filter).Decode(&user)
					fmt.Println("Found user: %v", user)
					if err != nil {
						log.Fatal(err)
					}

					return user, nil
				}
				return nil, nil
			},
		},
	}
	return fields
}
