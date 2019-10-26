package routes

import (
	"bytes"
	"context"
	"munchserver/secrets"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(secrets.GetMongoURI()))
	if err != nil {
		panic(err)
	}

	Db = client.Database(secrets.GetTestMongoDBName())
	// Inject db to tests
	tests.Db = Db

	code := m.Run()

	tests.ClearDB()
	os.Exit(code)
}

func TestInvalidLogin(t *testing.T) {
	tests.ClearDB()
	loginBody := []byte(`{"email": "invalid@email.com", "password": "notMyPassword"}`)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLoginHandler)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("login with invalid credentials expected status code of %v, but got %v", http.StatusUnauthorized, rr.Code)
	}
}

func TestValidRegister(t *testing.T) {
	tests.ClearDB()
	registerBody := []byte(`{"firstName":"some", "lastName": "tester", "email": "tester@example.com", "password": "password123", "dateOfBirth": "1969-04-20T05:00:00.000Z"}`)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(registerBody))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostRegisterHandler)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("register with valid information expected status code of %v, but got %v", http.StatusOK, rr.Code)
	}
}
