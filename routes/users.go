package routes

import (
	"context"
	"encoding/json"
	"log"
	"munchserver/middleware"
	"munchserver/models"
	"munchserver/queries"
	"munchserver/secrets"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"
)

type loginRequest struct {
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

type registerRequest struct {
	NameFirst   *string    `json:"firstName"`
	NameLast    *string    `json:"lastName"`
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
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
	_, err = Db.Collection("users").InsertOne(context.TODO(), registeredUser)

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
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Find user in database
	var user models.JSONUser
	err = Db.Collection("users").FindOne(context.TODO(), queries.UserWithEmail(*login.Email)).Decode(&user)
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
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
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

	// Create the response
	resp, err := json.Marshal(loginResponse{
		Token: jwtString,
	})
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
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
	err := Db.Collection("users").FindOne(r.Context(), queries.WithID(userID), queries.OptionsWithProjection(queries.ProfileProjection())).Decode(&user)
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
	err := Db.Collection("users").FindOne(r.Context(), queries.WithID(userID), queries.OptionsWithProjection(queries.UserProjection())).Decode(&user)
	if err != nil {
		log.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
