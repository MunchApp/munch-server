package middleware

import (
	"context"
	"munchserver/secrets"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

type key string

const UserKey key = "user"

func AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		auth := strings.Split(r.Header.Get("Authorization"), " ")
		if len(auth) == 2 && auth[0] == "Bearer" {
			tokenString := auth[1]
			var claims jwt.StandardClaims
			_, err := jwt.ParseWithClaims(tokenString, &claims, secrets.GetJWTSecret)

			if err == nil && claims.Valid() == nil {
				ctx = context.WithValue(r.Context(), UserKey, claims.Subject)
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
