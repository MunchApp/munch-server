package routes

import (
	"encoding/json"
	"io/ioutil"
	"munchserver/models"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

// TODO: Add test for checking required fields when adding a food truck

func TestFoodTrucksGetEmpty(t *testing.T) {
	tests.ClearDB()
	req, _ := http.NewRequest("GET", "/foodtrucks", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting all food trucks expected status code of %v, but got %v", expected, rr.Code)
	}
	body, _ := ioutil.ReadAll(rr.Body)
	if string(body) != "[]\n" {
		t.Errorf("expected empty array, but got %v", string(body))
	}
}

func TestFoodTrucksGetSingle(t *testing.T) {
	tests.ClearDB()
	foodTruck := models.JSONFoodTruck{
		ID: "test",
	}
	tests.AddFoodTruck(foodTruck)

	req, _ := http.NewRequest("GET", "/foodtrucks", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting all food trucks expected status code of %v, but got %v", expected, rr.Code)
	}

	var foodTrucks []models.JSONFoodTruck
	json.NewDecoder(rr.Body).Decode(&foodTrucks)
	if len(foodTrucks) != 1 {
		t.Errorf("expected array with one element, but got %v", foodTrucks)
	}
}

func TestFoodTruckGetInvalid(t *testing.T) {
	tests.ClearDB()
	req, _ := http.NewRequest("GET", "/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "invalid-id",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusNotFound
	if rr.Code != expected {
		t.Errorf("getting single invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTruckGetValid(t *testing.T) {
	tests.ClearDB()
	addFoodTruck := models.JSONFoodTruck{
		ID: "test",
	}
	tests.AddFoodTruck(addFoodTruck)
	req, _ := http.NewRequest("GET", "/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "test",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting single valid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
	var foodTruck models.JSONFoodTruck
	json.NewDecoder(rr.Body).Decode(&foodTruck)
	if foodTruck.ID != "test" {
		t.Errorf("expected food truck with id test, but got %v", foodTruck)
	}
}
