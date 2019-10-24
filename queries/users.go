package queries

import (
	"go.mongodb.org/mongo-driver/bson"
)

func UserWithEmail(email *string) bson.M {
	return bson.M{"email": email}
}
