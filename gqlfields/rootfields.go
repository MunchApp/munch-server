package gqlfields

import (
	"context"
	"log"
	"munchserver/models"
	"munchserver/queries"

	"github.com/graphql-go/graphql"
)

// Sub-queries underneath root query
var RootFields = graphql.Fields{
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
	// Get food truck info
	"foodtruck": &graphql.Field{
		Type: foodTruckType,
		Args: graphql.FieldConfigArgument{ // Can define specific foodtruck from here by name or address
			"name":    &graphql.ArgumentConfig{Type: graphql.String},
			"address": &graphql.ArgumentConfig{Type: graphql.String},
		},
		//logic, how to actually go and fetch data to return from db
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			name := p.Args["name"].(string)
			address := p.Args["address"].(string)
			var fTruck models.JSONFoodTruck

			// Get specific FoodTruck from db
			//from name
			err := Db.Collection("foodtrucks").FindOne(context.TODO(), queries.FoodTruckWithName(name)).Decode(&fTruck)
			if err != nil {
				log.Fatal(err)
			}

			//from address
			err2 := Db.Collection("foodtrucks").FindOne(context.TODO(), queries.FoodTruckWithAddress(address)).Decode(&fTruck)
			if err2 != nil {
				log.Fatal(err)
			}

			return fTruck, nil
		},
	},
}
