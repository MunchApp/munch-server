package queries

import (
	"go.mongodb.org/mongo-driver/bson"
)

func UserWithEmail(email string) bson.M {
	return bson.M{"email": email}
}

func PushOwnedFoodTruck(foodTruckID string) bson.M {
	return bson.M{"$push": bson.M{"ownedFoodTrucks": foodTruckID}}
}
