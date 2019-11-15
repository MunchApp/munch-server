package dbutils

import (
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateFoodTruck(avgRating float32, reviewID string) bson.M {
	return bson.M{"$set": bson.M{"avgRating": avgRating}, "$push": bson.M{"reviews": reviewID}}
}

func PushOwnedFoodTruck(foodTruckID string) bson.M {
	return bson.M{"$push": bson.M{"ownedFoodTrucks": foodTruckID}}
}

func PushReview(reviewID string) bson.M {
	return bson.M{"$push": bson.M{"reviews": reviewID}}
}
