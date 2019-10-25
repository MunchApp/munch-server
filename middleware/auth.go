package middleware

import (
	"context"
	"munchserver/secrets"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type key string

// UserKey is the key in request context of user's uuid
const UserKey key = "user"

// AuthenticateUser is a middleware which adds the authenticated user's uuid to the context of the request
func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get auth information from header
		auth := strings.Split(r.Header.Get("Authorization"), " ")
		if len(auth) == 2 && auth[0] == "Bearer" {
			tokenString := auth[1]

			// Get claims from jwt
			var claims jwt.StandardClaims
			_, err := jwt.ParseWithClaims(tokenString, &claims, secrets.GetJWTSecret)

			// If the token is still valid, add user to context
			if err == nil && claims.Valid() == nil {
				ctx = context.WithValue(r.Context(), UserKey, claims.Subject)
			}
		}

		// Go to next handler with new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
