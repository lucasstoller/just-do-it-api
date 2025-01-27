package routes

import (
	"just-do-it-api/handlers"
	"net/http"
)

func RegisterAuthRoutes() {
	http.HandleFunc("/api/auth/register", handlers.Register)
	http.HandleFunc("/api/auth/login", handlers.Login)
}
