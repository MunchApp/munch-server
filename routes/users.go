package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"munchserver/models"
	"munchserver/queries"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

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
		fmt.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create and insert user into database
	registeredUser := models.JSONUser{
		ID:           uuid.String(),
		PasswordHash: hashedPassword,
		NameFirst:    *newUser.NameFirst,
		NameLast:     *newUser.NameLast,
		Email:        *newUser.Email,
		DateOfBirth:  *newUser.DateOfBirth,
	}
	_, err = Db.Collection("users").InsertOne(context.TODO(), registeredUser)

	// If there is an error, it is most likely a duplicate user (email must be unique)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
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
	userResult := Db.Collection("users").FindOne(context.TODO(), queries.UserWithEmail(login.Email))
	err = userResult.Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if password matches
	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(*login.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a JWT for the user that expires in 15 minutes
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtString, err := token.SignedString([]byte("MunchIsReallyCool")) // TODO: Move the secret to an env var
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Create the response
	resp, err := json.Marshal(loginResponse{
		Token: jwtString,
	})
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
