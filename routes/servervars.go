package routes

import (
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Db     *mongo.Database
	Router *mux.Router
)
