package routes

import (
	"encoding/json"
	"log"
	"munchserver/middleware"
	"munchserver/models"
	"munchserver/queries"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
)

type addFoodTruckRequest struct {
	Name        *string       `json:"name"`
	Address     *string       `json:"address"`
	Location    [2]float64    `json:"location"`
	Hours       *[7][2]string `json:"hours"`
	Photos      *[]string     `json:"photos"`
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
		newFoodTruck.Hours == nil ||
		newFoodTruck.Photos == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate hours
	for i := 0; i < 7; i++ {
		validOpenTime, err := regexp.MatchString(`^\d{2}:\d{2}$`, newFoodTruck.Hours[i][0])
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		validCloseTime, err := regexp.MatchString(`^\d{2}:\d{2}$`, newFoodTruck.Hours[i][1])
		if err != nil {
			log.Printf("ERROR: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !validOpenTime || !validCloseTime {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// Generate uuid for food truck
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set tags to an empty array if they don't exist
	tags := newFoodTruck.Tags
	if tags == nil {
		tags = []string{}
	}

	addedFoodTruck := models.JSONFoodTruck{
		ID:          uuid.String(),
		Name:        *newFoodTruck.Name,
		Address:     *newFoodTruck.Address,
		Location:    newFoodTruck.Location,
		Owner:       user,
		Hours:       *newFoodTruck.Hours,
		Reviews:     []string{},
		Photos:      *newFoodTruck.Photos,
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
		_, err = Db.Collection("users").UpdateOne(r.Context(), queries.WithID(user), queries.PushOwnedFoodTruck(uuid.String()))
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
	err := Db.Collection("foodTrucks").FindOne(r.Context(), queries.WithID(foodTruckID)).Decode(&foodTruck)
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

	// Create correct filter
	filter := bson.D{}

	nameParam := r.URL.Query().Get("name")
	if nameParam != "" {
		filter = append(filter, bson.E{"name", nameParam})
	}

	tagsParam := r.URL.Query()["tags"]
	if len(tagsParam) != 0 {
		var tagsFilter bson.D
		tagsFilter = append(tagsFilter, bson.E{"$in", tagsParam})
		filter = append(filter, bson.E{"tags", tagsFilter})
	}

	addressParam := r.URL.Query().Get("address")
	if addressParam != "" {
		filter = append(filter, bson.E{"address", addressParam})
	}

	cur, err := foodTrucksCollection.Find(r.Context(), filter)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Get users from cursor, convert to empty slice if no users in DB
	var foodTrucks []models.JSONFoodTruck
	cur.All(r.Context(), &foodTrucks)
	if foodTrucks == nil {
		foodTrucks = make([]models.JSONFoodTruck, 0)
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

	// // Get user from context
	// _, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// // Check for a user
	// if !userLoggedIn {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	return
	// }

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

	_, err = Db.Collection("foodTrucks").UpdateOne(r.Context(), queries.WithID(foodTruckID), update)
	if err != nil {
		log.Fatal(err)
	}

	// Send response
	w.WriteHeader(http.StatusOK)

}
