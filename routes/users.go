package routes

import (
	"encoding/json"
	"log"
	"munchserver/dbutils"
	"munchserver/middleware"
	"munchserver/models"
	"munchserver/secrets"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type loginResponse struct {
	Token string          `json:"token"`
	User  models.JSONUser `json:"userObject"`
}

type registerRequest struct {
	NameFirst   *string    `json:"firstName"`
	NameLast    *string    `json:"lastName"`
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
}

type updateUserRequest struct {
	NameFirst   *string    `json:"firstName"`
	NameLast    *string    `json:"lastName"`
	PhoneNumber *string    `json:"phoneNumber"`
	City        *string    `json:"city"`
	State       *string    `json:"state"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
}

// PostRegisterHandler handles the logic for registering a user
func PostRegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Decode registered user's data
	userDecoder := json.NewDecoder(r.Body)
	userDecoder.DisallowUnknownFields()
	var newUser registerRequest
	err := userDecoder.Decode(&newUser)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure all fields in registered user are provided
	if newUser.NameFirst == nil ||
		newUser.NameLast == nil ||
		newUser.Email == nil ||
		newUser.Password == nil ||
		newUser.DateOfBirth == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Salt and hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*newUser.Password), bcrypt.DefaultCost)

	// Generate a random uuidv4
	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create and insert user into database
	registeredUser := models.JSONUser{
		ID:              uuid.String(),
		PasswordHash:    hashedPassword,
		NameFirst:       *newUser.NameFirst,
		NameLast:        *newUser.NameLast,
		Email:           *newUser.Email,
		DateOfBirth:     *newUser.DateOfBirth,
		Favorites:       []string{},
		Reviews:         []string{},
		OwnedFoodTrucks: []string{},
	}
	_, err = Db.Collection("users").InsertOne(r.Context(), registeredUser)

	// If there is an error, it is most likely a duplicate user (email must be unique)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Create the response
	w.WriteHeader(http.StatusOK)
}

// PostLoginHandler handles the logic for logging in
func PostLoginHandler(w http.ResponseWriter, r *http.Request) {
	// Decode login user
	userDecoder := json.NewDecoder(r.Body)
	userDecoder.DisallowUnknownFields()
	var login loginRequest
	err := userDecoder.Decode(&login)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Make sure all fields in request are provided
	if login.Email == nil ||
		login.Password == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Find user in database
	var user models.JSONUser
	err = Db.Collection("users").FindOne(r.Context(), dbutils.WithEmailQuery(*login.Email)).Decode(&user)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if password matches
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(*login.Password))
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a JWT for the user that expires in 15 minutes
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60).Unix(),
		Subject:   user.ID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, _ := secrets.GetJWTSecret(nil)
	jwtString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse{
		Token: jwtString,
		User:  user,
	})
}

func PutFavoriteHandler(w http.ResponseWriter, r *http.Request) {

	// Checks for food truck ID
	params := mux.Vars(r)
	foodTruckID, foodTruckIDExists := params["foodTruckID"]
	if !foodTruckIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get user from context
	userID, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user, or if the user agent is from the scraper
	if !userLoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var updateOperator string

	// Determine if an add or delete
	action := r.URL.Query().Get("action")
	if action == "add" {
		updateOperator = "$addToSet"
	} else if action == "delete" {
		updateOperator = "$pull"
	} else {
		log.Printf("ERROR: Incorrect action to edit user favorites.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	updateFavesFilter := bson.M{updateOperator: bson.M{"favorites": foodTruckID}}

	_, err := Db.Collection("users").UpdateOne(r.Context(), dbutils.WithIDQuery(userID), updateFavesFilter)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)

}

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	// Get user from context
	userID, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user, or if the user agent is from the scraper
	if !userLoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get user from database
	var user models.JSONUser
	err := Db.Collection("users").FindOne(r.Context(), dbutils.WithIDQuery(userID), dbutils.OptionsWithProjection(dbutils.ProfileProjection())).Decode(&user)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	// Get user id from route params
	params := mux.Vars(r)
	userID, userIDExists := params["userID"]
	if !userIDExists {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get user from database
	var user models.JSONUser
	err := Db.Collection("users").FindOne(r.Context(), dbutils.WithIDQuery(userID), dbutils.OptionsWithProjection(dbutils.UserProjection())).Decode(&user)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func PutUpdateUserHandler(w http.ResponseWriter, r *http.Request) {

	// Get user from context
	userID, userLoggedIn := r.Context().Value(middleware.UserKey).(string)

	// Check for a user
	if !userLoggedIn {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userDecoder := json.NewDecoder(r.Body)
	userDecoder.DisallowUnknownFields()

	// Decode request
	var updatedUser updateUserRequest
	err := userDecoder.Decode(&updatedUser)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Determine which fields should be updated
	var updateData bson.D

	if updatedUser.NameFirst != nil {
		updateData = append(updateData, bson.E{"firstName", *updatedUser.NameFirst})
	}
	if updatedUser.NameLast != nil {
		updateData = append(updateData, bson.E{"lastName", *updatedUser.NameLast})
	}
	if updatedUser.PhoneNumber != nil {
		updateData = append(updateData, bson.E{"phoneNumber", *updatedUser.PhoneNumber})
	}
	if updatedUser.City != nil {
		updateData = append(updateData, bson.E{"city", *updatedUser.City})
	}
	if updatedUser.State != nil {
		updateData = append(updateData, bson.E{"state", *updatedUser.State})
	}
	if updatedUser.DateOfBirth != nil {
		updateData = append(updateData, bson.E{"dateOfBirth", *updatedUser.DateOfBirth})
	}

	// Update food truck document
	update := bson.D{
		{"$set", updateData},
	}

	_, err = Db.Collection("users").UpdateOne(r.Context(), queries.WithID(userID), update)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Send response
	w.WriteHeader(http.StatusOK)

}
