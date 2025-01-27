package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"just-do-it-api/handlers"
	"just-do-it-api/middleware"
	"just-do-it-api/models"
)

func RegisterTaskRoutes() {
	// Base tasks endpoints
	http.HandleFunc("/v1/tasks", middleware.Logger(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetTasks(w, r)
		case http.MethodPost:
			handlers.CreateTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Method not allowed",
				"Method not supported for this endpoint",
			))
		}
	})))

	// Task operations by ID
	http.HandleFunc("/v1/tasks/", middleware.Logger(middleware.AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/v1/tasks/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Handle toggle completion endpoint
		if strings.HasSuffix(r.URL.Path, "/toggle") {
			if r.Method == http.MethodPatch {
				handlers.ToggleTask(w, r)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Handle regular CRUD operations
		switch r.Method {
		case http.MethodPut:
			handlers.UpdateTask(w, r)
		case http.MethodDelete:
			handlers.DeleteTask(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Method not allowed",
				"Method not supported for this endpoint",
			))
		}
	})))

	// Task filter endpoints
	http.HandleFunc("/v1/tasks/today", middleware.Logger(middleware.AuthMiddleware(handlers.GetTodayTasks)))
	http.HandleFunc("/v1/tasks/backlog", middleware.Logger(middleware.AuthMiddleware(handlers.GetBacklogTasks)))
}
