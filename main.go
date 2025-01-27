package main

import (
	"flag"
	"log"
	"net/http"

	"just-do-it-api/database"
	"just-do-it-api/routes"
)

func main() {
	// Parse command line flags
	reset := flag.Bool("reset", false, "Reset database and rerun all migrations")
	flag.Parse()

	// Handle database migrations
	if *reset {
		if err := database.ResetDatabase(""); err != nil {
			log.Fatalf("Failed to reset database: %v", err)
		}
	} else {
		if err := database.RunMigrations(""); err != nil {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	}

	// Initialize database connection
	database.CreateConnection()

	// Register routes
	routes.RegisterTaskRoutes()
	routes.RegisterAuthRoutes()

	log.Printf("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
