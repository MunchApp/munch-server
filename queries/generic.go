package queries

import "go.mongodb.org/mongo-driver/bson"

func All() bson.M {
	return bson.M{}
}

func WithID(id string) bson.M {
	return bson.M{"_id": id}
}

func WithFoodTruck(id string) bson.M {
	return bson.M{"foodTruck": id}
}
