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
	w.Write([]byte(addedFoodTruck.ID))
}

func GetFoodTrucksHandler(w http.ResponseWriter, r *http.Request) {
	// Get all foodtrucks from the database into a cursor
	foodTrucksCollection := Db.Collection("foodTrucks")
	cur, err := foodTrucksCollection.Find(r.Context(), bson.D{})
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

	// Convert users to json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(foodTrucks)
}
