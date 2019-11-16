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
	if foodTrucks[0].ID != "test" {
		t.Errorf("expected food truck to have id test, but got %v", foodTrucks[0].ID)
	}
}

func TestFoodTrucksGetInvalidLat(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/foodtrucks?lat=joe's house&lon=-97.735592", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting all food trucks with invalid latitude expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksGetInvalidLon(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/foodtrucks?lat=30.288441&lon=joe's house", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting all food trucks with invalid longitude expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksGetInvalidLatLon(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/foodtrucks?lat=-97.735592&lon=30.288441", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting all food trucks with invalid location expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksGetDistance(t *testing.T) {
	tests.ClearDB()
	foodTruck := models.JSONFoodTruck{
		ID:       "test",
		Location: [2]float64{-97.739928, 30.290241},
	}
	tests.AddFoodTruck(foodTruck)

	req, _ := http.NewRequest("GET", "/foodtrucks?lat=30.288441&lon=-97.735592", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting all food trucks expected status code of %v, but got %v", expected, rr.Code)
	}

	var foodTrucks []foodTruckWithDistance
	json.NewDecoder(rr.Body).Decode(&foodTrucks)
	if len(foodTrucks) != 1 {
		t.Errorf("expected array with one element, but got %v", foodTrucks)
	}
	if foodTrucks[0].ID != "test" {
		t.Errorf("expected food truck to have id test, but got %v", foodTrucks[0].ID)
	}
	if foodTrucks[0].Distance-462 > 1 {
		t.Errorf("expected food truck to have distance of around %v, but got %v", 462, foodTrucks[0].Distance)
	}
}

func TestFoodTrucksGetSortOrder(t *testing.T) {
	tests.ClearDB()
	foodTruck1 := models.JSONFoodTruck{
		ID:       "test1",
		Name:     "coop",
		Location: [2]float64{-97.742496, 30.286302},
	}
	foodTruck2 := models.JSONFoodTruck{
		ID:       "test2",
		Name:     "26th",
		Location: [2]float64{-97.744605, 30.290466},
	}
	foodTruck3 := models.JSONFoodTruck{
		ID:       "test3",
		Name:     "kins",
		Location: [2]float64{-97.739928, 30.290241},
	}
	tests.AddFoodTruck(foodTruck1)
	tests.AddFoodTruck(foodTruck2)
	tests.AddFoodTruck(foodTruck3)

	req, _ := http.NewRequest("GET", "/foodtrucks?lat=30.288441&lon=-97.735592", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting all food trucks expected status code of %v, but got %v", expected, rr.Code)
	}

	var foodTrucks []models.JSONFoodTruck
	json.NewDecoder(rr.Body).Decode(&foodTrucks)
	if len(foodTrucks) != 3 {
		t.Errorf("expected array with three element, but got %v", foodTrucks)
	}
	if foodTrucks[0].ID != "test3" {
		t.Errorf("expected first food truck to be kins, but got %v", foodTrucks[0].Name)
	}
	if foodTrucks[1].ID != "test1" {
		t.Errorf("expected second food truck to be coop, but got %v", foodTrucks[1].Name)
	}
	if foodTrucks[2].ID != "test2" {
		t.Errorf("expected third food truck to be 26th, but got %v", foodTrucks[2].Name)
	}
}

func TestFoodTruckGetEmpty(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/foodtruck", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting no food truck expected status code of %v, but got %v", expected, rr.Code)
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

func TestFoodTrucksGetSearchValid(t *testing.T) {
	tests.ClearDB()

	addFoodTruck := models.JSONFoodTruck{
		ID:   "test",
		Name: "testTruck",
	}
	tests.AddFoodTruck(addFoodTruck)

	req, _ := http.NewRequest("GET", "/foodtrucks?query=testTruck", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting valid food truck with name expected status code of %v, but got %v", expected, rr.Code)
	}

	var foodTruck []models.JSONFoodTruck
	json.NewDecoder(rr.Body).Decode(&foodTruck)

	if foodTruck[0].Name != "testTruck" {
		t.Errorf("expected food truck with name testTruck, but got %v", foodTruck[0].Name)
	}
}

func TestFoodTrucksGetSearchMultipleValid(t *testing.T) {
	tests.ClearDB()

	addFoodTruck1 := models.JSONFoodTruck{
		ID:   "test",
		Name: "ice cold",
	}
	addFoodTruck2 := models.JSONFoodTruck{
		ID:   "test",
		Name: "ice cream",
	}
	tests.AddFoodTruck(addFoodTruck1)
	tests.AddFoodTruck(addFoodTruck2)

	req, _ := http.NewRequest("GET", "/foodtrucks?query=ice+cream", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting valid food truck with name expected status code of %v, but got %v", expected, rr.Code)
	}

	var foodTrucks []models.JSONFoodTruck
	json.NewDecoder(rr.Body).Decode(&foodTrucks)

	if len(foodTrucks) != 1 {
		t.Errorf("expected one result from search of ice cream, but got %v", len(foodTrucks))
	}

	if foodTrucks[0].ID != "test" {
		t.Errorf("expected food truck with name ice cream, but got %v", foodTrucks[0].Name)
	}
}

func TestFoodTrucksPostValid(t *testing.T) {
	tests.ClearDB()

	tests.AddUser(models.JSONUser{
		ID:              "testuser",
		OwnedFoodTrucks: []string{},
	})

	name := "Luke's Coffee House"
	address := "2502 Nueces St\nAustin, TX 78705"
	location := [2]float64{-97.74731, 30.28793}
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
	phone := "8006729102"
	description := "testDescription"
	tags := []string{"food", "good food"}

	newFoodTruckTest := addFoodTruckRequest{
		Name:        &name,
		Address:     &address,
		Location:    &location,
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
		t.Errorf("adding valid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	var addedFoodTruck models.JSONReview
	json.NewDecoder(rr.Body).Decode(&addedFoodTruck)

	updatedFoodTruck := tests.GetFoodTruck(addedFoodTruck.ID)
	if updatedFoodTruck == nil || updatedFoodTruck.Name != "Luke's Coffee House" {
		t.Error("Error finding the added food truck in the database.")
	}
}

func TestFoodTrucksPostInvalidHours(t *testing.T) {
	tests.ClearDB()

	name := "Luke's Coffee House"
	address := "2502 Nueces St\nAustin, TX 78705"
	location := [2]float64{-97.74731, 30.28793}
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
	phone := "8006729102"
	description := "testDescription"
	tags := []string{"food", "good food"}

	newFoodTruckTest := addFoodTruckRequest{
		Name:        &name,
		Address:     &address,
		Location:    &location,
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
		t.Errorf("adding food truck with invalid hours expected status code of %v, but got %v", expected, rr.Code)
	}

}

func TestFoodTrucksPostUnauthorized(t *testing.T) {
	tests.ClearDB()

	foodTruckRequest := addFoodTruckRequest{}
	body, _ := json.Marshal(foodTruckRequest)

	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostFoodTrucksHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("adding food truck while unauthorized expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksPostInvalidBody(t *testing.T) {
	tests.ClearDB()

	foodTruckRequest := invalidRequestBody{}
	body, _ := json.Marshal(foodTruckRequest)

	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostFoodTrucksHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksPostInvalid(t *testing.T) {
	tests.ClearDB()

	foodTruckRequest := addFoodTruckRequest{}
	body, _ := json.Marshal(foodTruckRequest)

	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostFoodTrucksHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksPostNoPhotos(t *testing.T) {
	tests.ClearDB()

	name := "Luke's Coffee House"
	address := "2502 Nueces St\nAustin, TX 78705"
	location := [2]float64{-97.74731, 30.28793}
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
	phone := "8006729102"
	description := "testDescription"
	tags := []string{"food", "good food"}

	newFoodTruckTest := addFoodTruckRequest{
		Name:        &name,
		Address:     &address,
		Location:    &location,
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
		t.Errorf("adding food truck with no photos expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestFoodTrucksPostValidMinimum(t *testing.T) {
	tests.ClearDB()

	name := "test cafe"
	address := "1 direction st"
	location := [2]float64{0, 0}
	hours := [7][2]string{[2]string{"10:00", "10:00"}, [2]string{"10:00", "10:00"}, [2]string{"10:00", "10:00"}, [2]string{"10:00", "10:00"}, [2]string{"10:00", "10:00"}, [2]string{"10:00", "10:00"}, [2]string{"10:00", "10:00"}}
	photos := []string{"fake-url.com/image.png"}
	foodTruckRequest := addFoodTruckRequest{
		Name:     &name,
		Address:  &address,
		Location: &location,
		Hours:    &hours,
		Photos:   &photos,
	}
	body, _ := json.Marshal(foodTruckRequest)

	req, _ := http.NewRequest("POST", "/foodtrucks", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostFoodTrucksHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("adding valid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
	var foodTruck models.JSONFoodTruck
	json.NewDecoder(rr.Body).Decode(&foodTruck)

	addedFoodTruck := tests.GetFoodTruck(foodTruck.ID)
	if addedFoodTruck == nil || addedFoodTruck.Name != name || addedFoodTruck.Address != address {
		t.Error("adding valid food truck expected food truck in db")
	}
}

func TestClaimFoodTruckPutUnauthorized(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("PUT", "/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "test",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PutClaimFoodTruckHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("claiming a food truck while unauthorized expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestClaimFoodTruckPutNoFoodTruck(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("PUT", "/foodtruck", nil)
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PutClaimFoodTruckHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("claiming an invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestClaimFoodTruckPutInvalidFoodTruck(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("PUT", "/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "test",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PutClaimFoodTruckHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusNotFound
	if rr.Code != expected {
		t.Errorf("claiming an invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestClaimFoodTruckPutValid(t *testing.T) {
	tests.ClearDB()

	tests.AddUser(models.JSONUser{
		ID:              "testuser",
		OwnedFoodTrucks: []string{},
	})
	tests.AddFoodTruck(models.JSONFoodTruck{
		ID: "testfoodtruck",
	})

	req, _ := http.NewRequest("PUT", "/foodtruck", nil)
	vars := map[string]string{
		"foodTruckID": "testfoodtruck",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PutClaimFoodTruckHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("claiming a valid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	foodTruck := tests.GetFoodTruck("testfoodtruck")
	if foodTruck == nil || foodTruck.Owner != "testuser" {
		t.Error("claiming a valid food truck should have updated owner of food truck")
	}

	user := tests.GetUser("testuser")
	if user == nil || len(user.OwnedFoodTrucks) != 1 || user.OwnedFoodTrucks[0] != "testfoodtruck" {
		t.Error("claiming a valid food truck should have added food truck to user")
	}
}
