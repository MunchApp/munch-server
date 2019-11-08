package routes

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"munchserver/models"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReviewsGetEmpty(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/reviews", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	var reviews []models.JSONReview
	json.NewDecoder(rr.Body).Decode(&reviews)

	if len(reviews) != 0 {
		t.Errorf("getting all reviews of empty db expected 0 elements, but got %v", len(reviews))
	}
}

func TestReviewsGetNonEmpty(t *testing.T) {
	tests.ClearDB()

	review := models.JSONReview{
		ID: "testreview",
	}
	tests.AddReview(review)

	req, _ := http.NewRequest("GET", "/reviews", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	var reviews []models.JSONReview
	json.NewDecoder(rr.Body).Decode(&reviews)

	if len(reviews) != 1 {
		t.Errorf("getting all reviews expected 1 element, but got %v", len(reviews))
	}

	if reviews[0].ID != "testreview" {
		t.Errorf("expected review's id to be testreview but got %v", reviews[0].ID)
	}
}

func TestReviewsOfFoodTruckGetInvalidID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reviews/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "invalid-id",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsOfFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusNotFound
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsOfFoodTruckGetInvalid(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reviews/foodtruck", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsOfFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsOfFoodTruckGetValidWithReview(t *testing.T) {
	tests.ClearDB()
	testFoodTruck := models.JSONFoodTruck{
		ID:      "test",
		Reviews: []string{"test"},
	}
	testReview := models.JSONReview{
		ID:        "test",
		FoodTruck: "test",
	}
	tests.AddFoodTruck(testFoodTruck)
	tests.AddReview(testReview)

	req, _ := http.NewRequest("GET", "/reviews/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "test",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsOfFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting reviews of valid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	var reviews []models.JSONReview
	json.NewDecoder(rr.Body).Decode(&reviews)

	if len(reviews) != 1 {
		t.Errorf("getting reviews of valid food truck expected 1 element, but got %v", len(reviews))
	}

	if reviews[0].ID != "test" {
		t.Errorf("expected review's id to be test but got %v", reviews[0].ID)
	}
}

func TestReviewsOfFoodTruckGetValidWithoutReview(t *testing.T) {
	tests.ClearDB()
	testFoodTruck := models.JSONFoodTruck{
		ID:      "test",
		Reviews: []string{},
	}
	tests.AddFoodTruck(testFoodTruck)

	req, _ := http.NewRequest("GET", "/reviews/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "test",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsOfFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting reviews of valid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	var reviews []models.JSONReview
	json.NewDecoder(rr.Body).Decode(&reviews)

	if len(reviews) != 0 {
		t.Errorf("getting reviews of valid food truck expected 0 elements, but got %v", len(reviews))
	}
}
