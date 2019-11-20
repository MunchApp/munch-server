package routes

import (
	"bytes"
	"encoding/json"
	"log"
	"munchserver/dbutils"
	"munchserver/middleware"
	"munchserver/models"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type addFoodTruckRequest struct {
	Name        *string       `json:"name"`
	Address     *string       `json:"address"`
	Location    *[2]float64   `json:"location"`
	Hours       *[7][2]string `json:"hours"`
	Photos      []string      `json:"photos"`
	Website     string        `json:"website"`
	PhoneNumber string        `json:"phoneNumber"`
	Description string        `json:"description"`
	Tags        []string      `json:"tags"`
}

type updateFoodTruckRequest struct {
	Name        *string       `json:"name"`
	Address     *string       `json:"address"`
	Location    *[2]float64   `json:"location"`
	Status      *bool         `json:"status"`
	Hours       *[7][2]string `json:"hours"`
	Photos      []string      `json:"photos"`
	Website     *string       `json:"website"`
	PhoneNumber *string       `json:"phoneNumber"`
	Description *string       `json:"description"`
	Tags        []string      `json:"tags"`
}

type foodTruckWithDistance struct {
	ID          string       `json:"id" bson:"_id"`
	Name        string       `json:"name" bson:"name"`
	Address     string       `json:"address" bson:"address"`
	Location    [2]float64   `json:"location" bson:"location"`
	Owner       string       `json:"owner" bson:"owner"`
	Status      bool         `json:"status" bson:"status"`
	AvgRating   float32      `json:"avgRating" bson:"avgRating"`
	Hours       [7][2]string `json:"hours" bson:"hours"`
	Reviews     []string     `json:"reviews" bson:"reviews"`
	Photos      []string     `json:"photos" bson:"photos"`
	Website     string       `json:"website" bson:"website"`
	PhoneNumber string       `json:"phoneNumber" bson:"phoneNumber"`
	Description string       `json:"description" bson:"description"`
	Tags        []string     `json:"tags" bson:"tags"`
	Distance    float64      `json:"distance" bson:"distance"`
}

func PostFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	user, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user, or if the user agent is from the scraper
	if !userLoggedIn && r.Header.Get("User-Agent") != "MunchCritic/1.0" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	foodTruckDecoder := json.NewDecoder(r.Body)
	foodTruckDecoder.DisallowUnknownFields()

	// Decode request
	var newFoodTruck addFoodTruckRequest
	err := foodTruckDecoder.Decode(&newFoodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure required fields are set
	if newFoodTruck.Name == nil ||
		newFoodTruck.Address == nil ||
		newFoodTruck.Location == nil ||
		newFoodTruck.Hours == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate hours
	for i := 0; i < 7; i++ {
		validOpenTime, _ := regexp.MatchString(`^\d{2}:\d{2}$`, newFoodTruck.Hours[i][0])
		validCloseTime, _ := regexp.MatchString(`^\d{2}:\d{2}$`, newFoodTruck.Hours[i][1])
		if !validOpenTime || !validCloseTime {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Generate uuid for food truck
	uuid, _ := uuid.NewRandom()

	// Set tags to an empty array if they don't exist
	tags := newFoodTruck.Tags
	if tags == nil {
		tags = []string{}
	}
	photos := newFoodTruck.Photos
	if photos == nil {
		photos = []string{}
	}

	addedFoodTruck := models.JSONFoodTruck{
		ID:          uuid.String(),
		Name:        *newFoodTruck.Name,
		Address:     *newFoodTruck.Address,
		Location:    *newFoodTruck.Location,
		Owner:       user,
		Hours:       *newFoodTruck.Hours,
		Reviews:     []string{},
		Photos:      photos,
		Website:     newFoodTruck.Website,
		PhoneNumber: newFoodTruck.PhoneNumber,
		Description: newFoodTruck.Description,
		Tags:        tags,
	}

	// Add food truck to database
	_, err = Db.Collection("foodTrucks").InsertOne(r.Context(), addedFoodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update user that owns food truck
	if user != "" {
		_, err = Db.Collection("users").UpdateOne(r.Context(), dbutils.WithIDQuery(user), dbutils.PushOwnedFoodTruck(uuid.String()))
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addedFoodTruck)
}

func GetFoodTruckHandler(w http.ResponseWriter, r *http.Request) {
	// Get food truck id from route params
	params := mux.Vars(r)
	foodTruckID, foodTruckIDExists := params["foodTruckID"]
	if !foodTruckIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get food truck from database
	var foodTruck models.JSONFoodTruck
	err := Db.Collection("foodTrucks").FindOne(r.Context(), dbutils.WithIDQuery(foodTruckID)).Decode(&foodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foodTruck)
}

func GetFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {

	// Get all foodtrucks from the database into a cursor
	foodTrucksCollection := Db.Collection("foodTrucks")

	// Parse location from query params
	var location []float64
	if r.URL.Query().Get("lon") != "" || r.URL.Query().Get("lat") != "" {
		// Get location from query params
		longitude, err := strconv.ParseFloat(r.URL.Query().Get("lon"), 64)
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		latitude, err := strconv.ParseFloat(r.URL.Query().Get("lat"), 64)
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		location = []float64{longitude, latitude}
	}

	// Create correct filter
	var filter bson.M

	query := r.URL.Query().Get("query")

	// Create correct filter if have tags or name
	if query != "" {

		// Create name regex
		queryParams := strings.Split(query, " ")

		// Set regex as case insensitive
		nameRegex := "(?i)"
		for i, query := range queryParams {
			nameRegex += "(" + query + ")"
			if i < len(queryParams)-1 {
				nameRegex += "|"
			}

		}

		// Filter for tags and name
		tagsParam := []interface{}{
			bson.M{"tags": bson.M{"$regex": nameRegex}},
			bson.M{"name": bson.M{"$regex": nameRegex}},
		}
		filter = bson.M{"$or": tagsParam}

	} else {
		filter = dbutils.AllQuery()
	}

	var cur *mongo.Cursor
	var err error
	if location != nil {
		geoStage := bson.D{
			{"$geoNear", bson.M{
				"near": bson.M{
					"type":        "Point",
					"coordinates": location,
				},
				"distanceField": "distance",
				"spherical":     true,
				"query":         filter,
			}},
		}
		cur, err = foodTrucksCollection.Aggregate(r.Context(), mongo.Pipeline{geoStage})
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		cur, err = foodTrucksCollection.Find(r.Context(), filter)
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Get users from cursor, convert to empty slice if no users in DB
	var foodTrucks []foodTruckWithDistance
	cur.All(r.Context(), &foodTrucks)
	if foodTrucks == nil {
		foodTrucks = make([]foodTruckWithDistance, 0)
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foodTrucks)
}

func PutFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {

	// Checks for food truck ID
	params := mux.Vars(r)
	foodTruckID, foodTruckIDExists := params["foodTruckID"]
	if !foodTruckIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get user from context
	_, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user
	if !userLoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	foodTruckDecoder := json.NewDecoder(r.Body)
	foodTruckDecoder.DisallowUnknownFields()

	// Decode request
	var currentFoodTruck updateFoodTruckRequest
	err := foodTruckDecoder.Decode(&currentFoodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Determine which fields should be updated
	var updateData bson.D

	if currentFoodTruck.Name != nil {
		updateData = append(updateData, bson.E{"name", *currentFoodTruck.Name})
	}
	if currentFoodTruck.Address != nil {
		updateData = append(updateData, bson.E{"address", *currentFoodTruck.Address})
	}
	if currentFoodTruck.Location != nil {
		updateData = append(updateData, bson.E{"location", *currentFoodTruck.Location})
	}
	if currentFoodTruck.Status != nil {
		updateData = append(updateData, bson.E{"status", *currentFoodTruck.Status})
	}
	// Validate hours if updating
	if currentFoodTruck.Hours != nil {
		for i := 0; i < 7; i++ {
			validOpenTime, err := regexp.MatchString(`^\d{2}:\d{2}$`, currentFoodTruck.Hours[i][0])
			if err != nil {
				log.Printf("ERROR: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			validCloseTime, err := regexp.MatchString(`^\d{2}:\d{2}$`, currentFoodTruck.Hours[i][1])
			if err != nil {
				log.Printf("ERROR: %v", err)
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if !validOpenTime || !validCloseTime {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		}
		updateData = append(updateData, bson.E{"hours", *currentFoodTruck.Hours})
	}
	if currentFoodTruck.Photos != nil {
		updateData = append(updateData, bson.E{"photos", currentFoodTruck.Photos})
	}
	if currentFoodTruck.Website != nil {
		updateData = append(updateData, bson.E{"website", *currentFoodTruck.Website})
	}
	if currentFoodTruck.PhoneNumber != nil {
		updateData = append(updateData, bson.E{"phoneNumber", *currentFoodTruck.PhoneNumber})
	}
	if currentFoodTruck.Description != nil {
		updateData = append(updateData, bson.E{"description", *currentFoodTruck.Description})
	}
	if currentFoodTruck.Tags != nil {
		updateData = append(updateData, bson.E{"tags", currentFoodTruck.Tags})
	}

	// Update food truck document
	update := bson.D{
		{"$set", updateData},
	}

	_, err = Db.Collection("foodTrucks").UpdateOne(r.Context(), dbutils.WithIDQuery(foodTruckID), update)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

func PutClaimFoodTruckHandler(w http.ResponseWriter, r *http.Request) {

	// Checks for food truck ID
	params := mux.Vars(r)
	foodTruckID, foodTruckIDExists := params["foodTruckID"]
	if !foodTruckIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get user from context
	userID, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user
	if !userLoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Lookup food truck in db
	var foodTruck models.JSONFoodTruck
	err := Db.Collection("foodTrucks").FindOne(r.Context(), dbutils.WithIDQuery(foodTruckID)).Decode(&foodTruck)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Set food truck owner
	_, err = Db.Collection("foodTrucks").UpdateOne(r.Context(), dbutils.WithIDQuery(foodTruckID), dbutils.SetFoodTruckOwner(userID))
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add food truck to user
	_, err = Db.Collection("users").UpdateOne(r.Context(), dbutils.WithIDQuery(userID), dbutils.PushOwnedFoodTruck(foodTruckID))
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)
}

func PutFoodTruckUploadHandler(w http.ResponseWriter, r *http.Request) {

	// Checks for food truck ID
	params := mux.Vars(r)
	foodTruckID, foodTruckIDExists := params["foodTruckID"]
	if !foodTruckIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get user from context
	_, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user, or if the user agent is from the scraper
	if !userLoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(1024000 * 4)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusRequestEntityTooLarge)
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if filepath.Ext(fileHeader.Filename) != ".jpg" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer file.Close()
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// Generate a random uuidv4
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create file name for image
	filename := uuid.String() + filepath.Ext(fileHeader.Filename)

	// Upload image to s3
	result, err := Uploader.UploadWithContext(r.Context(), &s3manager.UploadInput{
		Bucket: aws.String("munch-assets"),
		Key:    aws.String(filename),
		Body:   bytes.NewReader(buffer),
	})
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err = Db.Collection("foodTrucks").UpdateOne(r.Context(), dbutils.WithIDQuery(foodTruckID), dbutils.PushPhoto(result.Location))
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
