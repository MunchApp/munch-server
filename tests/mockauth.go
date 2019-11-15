package tests

import (
	"context"
	"munchserver/middleware"
	"net/http"
)

// AuthenticateMockUser is a middleware which adds a mock user's uuid to the context of the request
func AuthenticateMockUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), middleware.UserKey, "testuser")

		// Go to next handler with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
