package middleware

import (
	"encoding/json"
	"just-do-it-api/models"
	"net/http"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Unauthorized",
				"Missing authentication token",
			))
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Unauthorized",
				"Invalid authentication token format",
			))
			return
		}

		token := bearerToken[1]
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Unauthorized",
				"Invalid or missing authentication token",
			))
			return
		}

		// TODO: Validate token with your auth service
		// For now, we'll just check if it's not empty

		next.ServeHTTP(w, r)
	}
}
