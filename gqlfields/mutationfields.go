package gqlfields

import (
	"context"
	"fmt"
	"log"
	"munchserver/models"

	"github.com/graphql-go/graphql"
)

// Sub-queries underneath mutation query
var MutationFields = graphql.Fields{
	"create": &graphql.Field{
		Type: UserType,
		Args: graphql.FieldConfigArgument{ //Can define specific user from here
			"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			firstname, ok := p.Args["name"].(string)
			if ok {

				fmt.Println("--- in mutation ---")
				collection := Db.Collection("users")

				// Insert to MongoDB
				// Sample user
				user1 := models.JSONUser{
					ID:        "1",
					NameFirst: firstname,
					Email:     "devil@gmail.com",
				}
				result, err := collection.InsertOne(context.TODO(), user1)
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println("Inserted a single doc:", result.InsertedID)

				return user1, nil
			}
			return nil, nil
		},
	},
}
