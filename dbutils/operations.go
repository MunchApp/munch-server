package dbutils

import (
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateFoodTruckWithReview(avgRating float32, reviewID string) bson.M {
	return bson.M{"$set": bson.M{"avgRating": avgRating}, "$push": bson.M{"reviews": reviewID}}
}

func SetFoodTruckOwner(userID string) bson.M {
	return bson.M{"$set": bson.M{"owner": userID}}
}

func SetProfilePicture(pictureURL string) bson.M {
	return bson.M{"$set": bson.M{"picture": pictureURL}}
}

func PushOwnedFoodTruck(foodTruckID string) bson.M {
	return bson.M{"$push": bson.M{"ownedFoodTrucks": foodTruckID}}
}

func PushReview(reviewID string) bson.M {
	return bson.M{"$push": bson.M{"reviews": reviewID}}
}

func PushPhoto(photoURL string) bson.M {
	return bson.M{"$push": bson.M{"photos": photoURL}}
}
