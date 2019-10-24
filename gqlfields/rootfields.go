package gqlfields

import (
	"context"
	"log"
	"munchserver/models"
	"munchserver/queries"

	"github.com/graphql-go/graphql"
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
				email := p.Args["email"].(string)
				var user models.JSONUser

				// Get specific user from db
				err := Db.Collection("users").FindOne(context.TODO(), queries.UserWithEmail(email)).Decode(&user)
				if err != nil {
					log.Fatal(err)
				}

				return user, nil
			},
		},
	}
	return fields
}
