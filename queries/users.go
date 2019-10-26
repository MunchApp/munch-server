package queries

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UserWithEmail(email string) bson.M {
	return bson.M{"email": email}
}

func PushOwnedFoodTruck(foodTruckID string) bson.M {
	return bson.M{"$push": bson.M{"ownedFoodTrucks": foodTruckID}}
}

func PushReview(reviewID string) bson.M {
	return bson.M{"$push": bson.M{"reviews": reviewID}}
}

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
