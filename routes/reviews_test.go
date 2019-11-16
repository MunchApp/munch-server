package routes

import (
	"bytes"
	"encoding/json"
	"munchserver/models"
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
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
	tests.ClearDB()

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
	tests.ClearDB()

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

func TestReviewsPostUnauthorized(t *testing.T) {
	tests.ClearDB()

	reviewsRequest := newReviewRequest{}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostReviewsHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusUnauthorized
	if rr.Code != expected {
		t.Errorf("adding review while unauthroized expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsPostInvalidRequestClient(t *testing.T) {
	tests.ClearDB()

	reviewsRequest := newReviewRequest{
		Comment: "Amazing food",
	}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostReviewsHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding review with invalid body expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsPostInvalidRequestScraper(t *testing.T) {
	tests.ClearDB()

	reviewsRequest := newReviewRequest{
		Comment: "Amazing food",
	}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	req.Header.Set("User-Agent", "MunchCritic/1.0")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostReviewsHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding review as scraper with invalid body expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsPostInvalidBody(t *testing.T) {
	tests.ClearDB()

	reviewsRequest := invalidRequestBody{}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostReviewsHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("adding review with invalid body expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsPostInvalidFoodTruck(t *testing.T) {
	tests.ClearDB()

	var rating float32 = 5.0
	name := "invalid-truck"
	reviewsRequest := newReviewRequest{
		FoodTruck: &name,
		Comment:   "Amazing food",
		Rating:    &rating,
	}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostReviewsHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusNotFound
	if rr.Code != expected {
		t.Errorf("adding review with invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}
}

func TestReviewsPostValidClient(t *testing.T) {
	tests.ClearDB()

	tests.AddFoodTruck(models.JSONFoodTruck{
		ID:        "testfoodtruck",
		Reviews:   []string{},
		AvgRating: 0.0,
	})
	tests.AddUser(models.JSONUser{
		ID:      "testuser",
		Reviews: []string{},
	})

	var rating float32 = 5.0
	name := "testfoodtruck"
	reviewsRequest := newReviewRequest{
		FoodTruck: &name,
		Comment:   "Amazing food",
		Rating:    &rating,
	}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostReviewsHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("adding review with invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	updatedFoodTruck := tests.GetFoodTruck("testfoodtruck")
	if updatedFoodTruck == nil || len(updatedFoodTruck.Reviews) != 1 {
		t.Error("adding valid review did not update food truck")
	}
	if updatedFoodTruck.AvgRating != 5.0 {
		t.Error("adding valid review did not update rating of food truck")
	}

	updatedUser := tests.GetUser("testuser")
	if updatedUser == nil || len(updatedUser.Reviews) != 1 {
		t.Error("adding valid review did not update user")
	}

	var addedReturnedReview models.JSONReview
	json.NewDecoder(rr.Body).Decode(&addedReturnedReview)

	addedReview := tests.GetReview(addedReturnedReview.ID)
	if addedReview == nil {
		t.Error("adding valid review did not add review to db")
	}
}

func TestReviewsPostValidScraper(t *testing.T) {
	tests.ClearDB()

	tests.AddFoodTruck(models.JSONFoodTruck{
		ID:      "testfoodtruck",
		Reviews: []string{},
	})
	var rating float32 = 5.0
	name := "testfoodtruck"
	reviewsRequest := newReviewRequest{
		ReviewerName: "Test User",
		FoodTruck:    &name,
		Comment:      "Amazing food",
		Rating:       &rating,
		Origin:       "Yelp",
	}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	req.Header.Set("User-Agent", "MunchCritic/1.0")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostReviewsHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("adding review with invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	updatedFoodTruck := tests.GetFoodTruck("testfoodtruck")
	if updatedFoodTruck == nil || len(updatedFoodTruck.Reviews) != 1 {
		t.Error("adding valid review did not update food truck")
	}

	var addedReturnedReview models.JSONReview
	json.NewDecoder(rr.Body).Decode(&addedReturnedReview)

	addedReview := tests.GetReview(addedReturnedReview.ID)
	if addedReview == nil {
		t.Error("adding valid review did not add review to db")
	}
}

func TestReviewsPostNewRating(t *testing.T) {
	tests.ClearDB()

	tests.AddFoodTruck(models.JSONFoodTruck{
		ID:        "testfoodtruck",
		Reviews:   []string{"fakereview"},
		AvgRating: 4.0,
	})
	tests.AddUser(models.JSONUser{
		ID:      "testuser",
		Reviews: []string{},
	})

	var rating float32 = 5.0
	name := "testfoodtruck"
	reviewsRequest := newReviewRequest{
		FoodTruck: &name,
		Comment:   "Amazing food",
		Rating:    &rating,
	}
	body, _ := json.Marshal(reviewsRequest)

	req, _ := http.NewRequest("POST", "/reviews", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()
	handler := tests.AuthenticateMockUser(http.HandlerFunc(PostReviewsHandler))
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("adding review with invalid food truck expected status code of %v, but got %v", expected, rr.Code)
	}

	updatedFoodTruck := tests.GetFoodTruck("testfoodtruck")
	if updatedFoodTruck == nil || len(updatedFoodTruck.Reviews) != 2 {
		t.Error("adding valid review did not update food truck")
	}
	if updatedFoodTruck.AvgRating != 4.5 {
		t.Errorf("adding valid review did not update rating of food truck, expected %v but got %v", 4.5, updatedFoodTruck.AvgRating)
	}
}

func TestReviewGet(t *testing.T) {
	tests.ClearDB()

	// Add sample review to DB
	tests.AddReview(models.JSONReview{
		ID: "test",
	})
	req, _ := http.NewRequest("GET", "/review", nil)
	vars := map[string]string{
		"reviewID": "test",
	}
	req = mux.SetURLVars(req, vars)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting single review expected status code of %v, but got %v", expected, rr.Code)
	}

	var review models.JSONReview
	json.NewDecoder(rr.Body).Decode(&review)
	if review.ID != "test" {
		t.Errorf("expected review with id test, but got %v", review.ID)
	}

}

func TestReviewGetEmpty(t *testing.T) {
	tests.ClearDB()

	// Add sample review to DB
	tests.AddReview(models.JSONReview{
		ID: "test",
	})
	req, _ := http.NewRequest("GET", "/review", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetReviewHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusBadRequest
	if rr.Code != expected {
		t.Errorf("getting single review expected status code of %v, but got %v", expected, rr.Code)
	}

}
