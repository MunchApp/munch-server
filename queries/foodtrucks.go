package queries

import (
	"go.mongodb.org/mongo-driver/bson"
)

func FoodTruckWithName(name string) bson.M {
	return bson.M{"name": name}
}
func FoodTruckWithAddress(address string) bson.M {
	return bson.M{"address": address}
}
