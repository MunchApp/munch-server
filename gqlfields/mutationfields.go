package gqlfields

import (
	"context"
	"fmt"
	"log"
	"munchserver/models"
	"munchserver/routes"

	"github.com/graphql-go/graphql"
)

// Sub-queries underneath root query
// func GetMutationFields() graphql.Fields {
// 	fields := graphql.Fields{
// 		"create": &graphql.Field{
// 			Type: UserType,
// 			Args: graphql.FieldConfigArgument{ //Can define specific user from here
// 				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
// 			},
// 			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
// 				firstname, ok := p.Args["name"].(string)
// 				if ok {

// 					fmt.Println("--- in mutation ---")
// 					collection := routes.Db.Collection("users")

// 					// Insert to MongoDB
// 					// Sample user
// 					user1 := models.JSONUser{"1", firstname, nil, nil, "devil@gmail.com", nil}
// 					result, err2 := collection.InsertOne(context.TODO(), user1)
// 					if err2 != nil {
// 						log.Fatal("user1 input %v", err2)
// 					}
// 					fmt.Println("Inserted a single doc:", result.InsertedID)

// 					return user1, nil
// 				}
// 				return nil, nil
// 			},
// 		},
// 	}
// 	return fields
// }

var MutationTypes = graphql.Fields{
	"create": &graphql.Field{
		Type: UserType,
		Args: graphql.FieldConfigArgument{ //Can define specific user from here
			"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			firstname, ok := p.Args["name"].(string)
			if ok {

				fmt.Println("--- in mutation ---")
				collection := routes.Db.Collection("users")

				// Insert to MongoDB
				// Sample user
				user1 := models.JSONUser{"1", firstname, nil, nil, "devil@gmail.com", nil}
				result, err2 := collection.InsertOne(context.TODO(), user1)
				if err2 != nil {
					log.Fatal("user1 input %v", err2)
				}
				fmt.Println("Inserted a single doc:", result.InsertedID)

				return user1, nil
			}
			return nil, nil
		},
	},
}
