package routes

import (
	"munchserver/models"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestReviewsOfFoodTruckGetInvalidID(t *testing.T) {
	req, _ := http.NewRequest("GET", "/reviews/foodtruck/", nil)
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
	req, _ := http.NewRequest("GET", "/reviews/foodtruck/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewsOfFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting reviews of invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsOfFoodTruckGetValid(t *testing.T) {
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

	req, _ := http.NewRequest("GET", "/reviews/foodtruck/", nil)
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
}
