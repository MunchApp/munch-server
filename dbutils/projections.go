package dbutils

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProfileProjection() bson.M {
	return bson.M{
		"passwordHash": 0,
	}
}

func UserProjection() bson.M {
	return bson.M{
		"passwordHash": 0,
		"dateOfBirth":  0,
		"phoneNumber":  0,
	}
}

func OptionsWithProjection(proj bson.M) *options.FindOneOptions {
	return &options.FindOneOptions{Projection: proj}
}
