package dbutils

import "go.mongodb.org/mongo-driver/bson"

func AllQuery() bson.M {
	return bson.M{}
}

func WithIDQuery(id string) bson.M {
	return bson.M{"_id": id}
}

func WithIDsQuery(ids []string) bson.M {
	return bson.M{"_id": bson.M{"$in": ids}}
}

func WithFoodTruckQuery(id string) bson.M {
	return bson.M{"foodTruck": id}
}

func WithEmailQuery(email string) bson.M {
	return bson.M{"email": email}
}
