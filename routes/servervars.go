package routes

import (
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	Db       *mongo.Database
	Router   *mux.Router
	Uploader *s3manager.Uploader
)
