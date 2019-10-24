package routes

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"munchserver/models"
	"net/http"
	"time"

	"github.com/google/uuid"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

type registerUser struct {
	NameFirst   *string    `json:"firstName"`
	NameLast    *string    `json:"lastName"`
	Email       *string    `json:"email"`
	Password    *string    `json:"password"`
	DateOfBirth *time.Time `json:"dateOfBirth"`
}

func PostUsersHandler(w http.ResponseWriter, r *http.Request) {
	userDecoder := json.NewDecoder(r.Body)
	userDecoder.DisallowUnknownFields()
	var newUser registerUser
	err := userDecoder.Decode(&newUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
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

	uuid, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	registeredUser := models.JSONUser{
		ID:           uuid.String(),
		PasswordHash: hashedPassword,
		NameFirst:    *newUser.NameFirst,
		NameLast:     *newUser.NameLast,
		Email:        *newUser.Email,
		DateOfBirth:  *newUser.DateOfBirth,
	}
	_, err = Db.Collection("users").InsertOne(context.TODO(), registeredUser)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Get all users from the database into a cursor
	usersCollection := Db.Collection("users")
	cur, err := usersCollection.Find(context.TODO(), bson.D{})
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in database: %v", err)
		return
	}

	// Get users from cursor, convert to empty slice if no users in DB
	var users []models.JSONUser
	cur.All(context.TODO(), &users)
	if users == nil {
		users = make([]models.JSONUser, 0)
	}

	// Convert users to json
	js, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in decoding mongo document: %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
