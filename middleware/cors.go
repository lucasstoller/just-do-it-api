package middleware

import (
	"net/http"

	"github.com/rs/cors"
)

// CorsMiddleware returns a CORS middleware configured for development
func CorsMiddleware() func(http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                            // Allow all origins in development
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"}, // Allow common HTTP methods
		AllowedHeaders:   []string{"*"},                            // Allow all headers
		AllowCredentials: true,                                     // Allow credentials (cookies, authorization headers, etc)
		Debug:            true,                                     // Enable debugging for development
	}).Handler
}
