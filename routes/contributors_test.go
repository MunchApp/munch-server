package routes

import (
	"munchserver/tests"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestContributorsGet(t *testing.T) {
	tests.ClearDB()

	req, _ := http.NewRequest("GET", "/contributors", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetContributorsHandler)
	handler.ServeHTTP(rr, req)

	expected := http.StatusOK
	if rr.Code != expected {
		t.Errorf("getting contributors expected status code of %v, but got %v", expected, rr.Code)
	}
}
