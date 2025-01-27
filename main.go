package main

import (
	"log"
	"net/http"

	"just-do-it-api/database"
	"just-do-it-api/models"
	"just-do-it-api/routes"
)

func main() {
	db := database.CreateConnection()
	db.AutoMigrate(&models.Task{})

	// Register routes
	routes.RegisterTaskRoutes()

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
