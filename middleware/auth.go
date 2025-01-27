package middleware

import (
	"context"
	"encoding/json"
	"just-do-it-api/auth"
	"just-do-it-api/models"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "userID"

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

		claims, err := auth.ValidateToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Unauthorized",
				"Invalid token",
			))
			return
		}

		// Add user ID to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

// GetUserID retrieves the user ID from the request context
func GetUserID(r *http.Request) uint {
	userID, _ := r.Context().Value(UserIDKey).(uint)
	return userID
}
