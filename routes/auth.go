package routes

import (
	"just-do-it-api/handlers"
	"net/http"
)

func RegisterAuthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/auth/register", handlers.Register)
	mux.HandleFunc("/api/auth/login", handlers.Login)
}
