package routes

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"munchserver/models"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

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

func TestPostFoodTruckValid(t *testing.T) {
	tests.ClearDB()

	name := "Luke's Coffee House"
	address := "2502 Nueces St\nAustin, TX 78705"
	location := [2]float64{30.28793, -97.74731}
	hours := [7][2]string{
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
	}
	photos := []string{
		"https://s3-media3.fl.yelpcdn.com/bphoto/d-1vKOqcEcRutpQ--jPa9A/o.jpg",
		"https://s3-media4.fl.yelpcdn.com/bphoto/vXk0bXpV007fSWQiOTlHgg/o.jpg",
		"https://s3-media2.fl.yelpcdn.com/bphoto/tUJ5gLnfRFhp_v-LUGj8Ww/o.jpg",
	}
	website := "www.google.com"
	phone := "+18006729102"
	description := "testDescription"
	tags := []string{"food", "good food"}

	// type addFoodTruckRequest struct {
	// 	Name        *string       `json:"name"`
	// 	Address     *string       `json:"address"`
	// 	Location    [2]float64    `json:"location"`
	// 	Hours       *[7][2]string `json:"hours"`
	// 	Photos      *[]string     `json:"photos"`
	// 	Website     string        `json:"website"`
	// 	PhoneNumber string        `json:"phoneNumber"`
	// 	Description string        `json:"description"`
	// 	Tags        []string      `json:"tags"`
	// }

	newFoodTruckTest := addFoodTruckRequest{
		Name:        &name,
		Address:     &address,
		Location:    location,
		Hours:       &hours,
		Photos:      &photos,
		Website:     website,
		PhoneNumber: phone,
		Description: description,
		Tags:        tags,
	}

	body, _ := json.Marshal(newFoodTruckTest)
	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostFoodTrucksHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("adding food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	var addedFoodTruck models.JSONReview
	json.NewDecoder(rr.Body).Decode(&addedFoodTruck)

	updatedFoodTruck := tests.GetFoodTruck(addedFoodTruck.ID)
	if updatedFoodTruck == nil || updatedFoodTruck.Name != "Luke's Coffee House" {
		t.Error("Error finding the added food truck in the database.")
		print(updatedFoodTruck.Name)

		if updatedFoodTruck == nil {
			t.Error("food truck is null")
		}
	}
}

func TestPostFoodTruckInvalidHours(t *testing.T) {
	tests.ClearDB()

	name := "Luke's Coffee House"
	address := "2502 Nueces St\nAustin, TX 78705"
	location := [2]float64{30.28793, -97.74731}
	hours := [7][2]string{
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"cute_string", "canthisbeparsed"},
	}
	photos := []string{
		"https://s3-media3.fl.yelpcdn.com/bphoto/d-1vKOqcEcRutpQ--jPa9A/o.jpg",
		"https://s3-media4.fl.yelpcdn.com/bphoto/vXk0bXpV007fSWQiOTlHgg/o.jpg",
		"https://s3-media2.fl.yelpcdn.com/bphoto/tUJ5gLnfRFhp_v-LUGj8Ww/o.jpg",
	}
	website := "www.google.com"
	phone := "+18006729102"
	description := "testDescription"
	tags := []string{"food", "good food"}

	// type addFoodTruckRequest struct {
	// 	Name        *string       `json:"name"`
	// 	Address     *string       `json:"address"`
	// 	Location    [2]float64    `json:"location"`
	// 	Hours       *[7][2]string `json:"hours"`
	// 	Photos      *[]string     `json:"photos"`
	// 	Website     string        `json:"website"`
	// 	PhoneNumber string        `json:"phoneNumber"`
	// 	Description string        `json:"description"`
	// 	Tags        []string      `json:"tags"`
	// }

	newFoodTruckTest := addFoodTruckRequest{
		Name:        &name,
		Address:     &address,
		Location:    location,
		Hours:       &hours,
		Photos:      &photos,
		Website:     website,
		PhoneNumber: phone,
		Description: description,
		Tags:        tags,
	}

	body, _ := json.Marshal(newFoodTruckTest)
	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostFoodTrucksHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding food truck expected status code of %v, but got %v", expected, rr.Code)
	}

}

func TestPostFoodTruckNoPhotos(t *testing.T) {
	tests.ClearDB()

	name := "Luke's Coffee House"
	address := "2502 Nueces St\nAustin, TX 78705"
	location := [2]float64{30.28793, -97.74731}
	hours := [7][2]string{
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
		[2]string{"10:00", "11:00"},
	}
	website := "www.google.com"
	phone := "+18006729102"
	description := "testDescription"
	tags := []string{"food", "good food"}

	// type addFoodTruckRequest struct {
	// 	Name        *string       `json:"name"`
	// 	Address     *string       `json:"address"`
	// 	Location    [2]float64    `json:"location"`
	// 	Hours       *[7][2]string `json:"hours"`
	// 	Photos      *[]string     `json:"photos"`
	// 	Website     string        `json:"website"`
	// 	PhoneNumber string        `json:"phoneNumber"`
	// 	Description string        `json:"description"`
	// 	Tags        []string      `json:"tags"`
	// }

	newFoodTruckTest := addFoodTruckRequest{
		Name:        &name,
		Address:     &address,
		Location:    location,
		Hours:       &hours,
		Website:     website,
		PhoneNumber: phone,
		Description: description,
		Tags:        tags,
	}

	body, _ := json.Marshal(newFoodTruckTest)
	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostFoodTrucksHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding food truck expected status code of %v, but got %v", expected, rr.Code)
	}

}