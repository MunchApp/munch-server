package gqlfields

import "github.com/graphql-go/graphql"

// GQL Fields to RETURN to user
var UserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "userFields",
	Fields: graphql.Fields{
		"firstName":  &graphql.Field{Type: graphql.String},
		"lastName":  &graphql.Field{Type: graphql.String},
		"email": &graphql.Field{Type: graphql.String},
		"id":    &graphql.Field{Type: graphql.String},
	},
})

// GQL Fields to RETURN for foodTruck
var FoodTruckType = graphql.NewObject(graphql.ObjectConfig{
	Name: "foodTruckFields",
	Fields: graphql.Fields{
		"name":  &graphql.Field{Type: graphql.String},
		"address": &graphql.Field{Type: graphql.String},
		"avgRating": &graphql.Field{Type: graphql.Float},
		"hours": &graphql.Field{Type: graphql.NewList(graphql.String)},
		"reviews": &graphql.Field{Type: graphql.NewList(ReviewType)},
	},
})

// GQL Fields to RETURN for reviews
var ReviewType = graphql.NewObject(graphql.ObjectConfig{
	Name: "reviewFields",
	Fields: graphql.Fields{
		"id":  &graphql.Field{Type: graphql.String},
		"reviewer": &graphql.Field{Type: graphql.String},              //maybe change to user type? Have to see how its done again
		"comment": &graphql.Field{Type: graphql.String},
		"rating": &graphql.Field{Type: graphql.Float},
		"date": &graphql.Field{Type: graphql.String},
	},
})
