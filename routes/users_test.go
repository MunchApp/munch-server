package routes

import (
	"bytes"
	"encoding/json"
	"munchserver/middleware"
	"munchserver/models"
	"munchserver/secrets"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginPostInvalidCredentials(t *testing.T) {
	tests.ClearDB()

	// Create login request
	email := "invalid@example.com"
	password := "notMyPassword"
	loginBody := loginRequest{
		Email:    &email,
		Password: &password,
	}

	body, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLoginHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("login with invalid credentials expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestLoginPostInvalidRequest(t *testing.T) {
	tests.ClearDB()

	loginBody := loginRequest{}
	body, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLoginHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("login with invalid request expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestLoginPostValid(t *testing.T) {
	tests.ClearDB()

	// Add user to db
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	dob, _ := time.Parse(time.RFC3339, "1969-04-20T05:00:00.000Z")
	addUser := models.JSONUser{
		ID:           "testuser",
		PasswordHash: hashedPassword,
		NameFirst:    "test",
		NameLast:     "user",
		Email:        "tester@example.com",
		DateOfBirth:  dob,
	}
	tests.AddUser(addUser)

	// Create login request
	email := "tester@example.com"
	password := "password123"
	loginBody := loginRequest{
		Email:    &email,
		Password: &password,
	}

	body, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLoginHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("login with valid credentials expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestLoginPostIncorrectPassword(t *testing.T) {
	tests.ClearDB()

	// Add user to db
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	dob, _ := time.Parse(time.RFC3339, "1969-04-20T05:00:00.000Z")
	addUser := models.JSONUser{
		ID:           "testuser",
		PasswordHash: hashedPassword,
		NameFirst:    "test",
		NameLast:     "user",
		Email:        "tester@example.com",
		DateOfBirth:  dob,
	}
	tests.AddUser(addUser)

	// Create login request
	email := "tester@example.com"
	password := "password"
	loginBody := loginRequest{
		Email:    &email,
		Password: &password,
	}

	body, _ := json.Marshal(loginBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLoginHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("login with valid credentials expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestRegisterPostValid(t *testing.T) {
	tests.ClearDB()

	// Create request body
	name := "tester"
	email := "tester@example.com"
	password := "password123"
	dob, _ := time.Parse(time.RFC3339, "1969-04-20T05:00:00.000Z")
	registerBody := registerRequest{
		NameFirst:   &name,
		NameLast:    &name,
		Email:       &email,
		Password:    &password,
		DateOfBirth: &dob,
	}
	body, _ := json.Marshal(registerBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostRegisterHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("register with valid information expected status code of %v, but got %v", expected, rr.Code)
	}

	// TODO: Add check for id returned
}

func TestRegisterPostInvalid(t *testing.T) {
	tests.ClearDB()

	// Create request body
	registerBody := registerRequest{}
	body, _ := json.Marshal(registerBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostRegisterHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("register with valid information expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestRegisterPostDuplicate(t *testing.T) {
	tests.ClearDB()

	// Add user to db
	user := models.JSONUser{
		Email: "tester@example.com",
	}
	tests.AddUser(user)

	// Create request body
	name := "tester"
	email := "tester@example.com"
	password := "password123"
	dob, _ := time.Parse(time.RFC3339, "1969-04-20T05:00:00.000Z")
	registerBody := registerRequest{
		NameFirst:   &name,
		NameLast:    &name,
		Email:       &email,
		Password:    &password,
		DateOfBirth: &dob,
	}
	body, _ := json.Marshal(registerBody)

	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostRegisterHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusConflict
	if rr.Code != expected {
		t.Errorf("register with duplicate email expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestProfileGetAuthorized(t *testing.T) {
	tests.ClearDB()

	// Add user to db
	user := models.JSONUser{
		ID: "testuser",
	}
	tests.AddUser(user)

	// Create a JWT
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		Subject:   "testuser",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, _ := secrets.GetJWTSecret(nil)
	jwtString, _ := token.SignedString(jwtSecret)

	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+jwtString)

	rr := httptest.NewRecorder()
	handler := middleware.AuthenticateUser(http.HandlerFunc(GetProfileHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting profile while logged in expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestProfileGetUnauthorized(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/profile", nil)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := middleware.AuthenticateUser(http.HandlerFunc(GetProfileHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("getting profile while not logged in expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestUserGetInvalidID(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/users", nil)
	vars := map[string]string{
		"userID": "invalid-id",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusNotFound
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestUserGetNoID(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestUserGetValidID(t *testing.T) {
	tests.ClearDB()

	user := models.JSONUser{
		ID: "testuser",
	}
	tests.AddUser(user)

	req, _ := http.NewRequest("GET", "/users", nil)
	vars := map[string]string{
		"userID": "testuser",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetUserHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}
