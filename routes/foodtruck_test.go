package routes

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFoodTruckGet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/foodtrucks", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetFoodTrucksHandler)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Errorf("getting all food trucks expected status code of %v, but got %v", http.StatusOK, rr.Code)
	}
	body, _ := ioutil.ReadAll(rr.Body)
	if string(body) != "[]\n" {
		t.Errorf("expected empty array, but got %v", string(body))
	}
}
