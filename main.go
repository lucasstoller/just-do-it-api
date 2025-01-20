package main

import (
	"just-do-it-api/database"
	"just-do-it-api/handlers"
	"just-do-it-api/models"
	"log"
	"net/http"
)

func main() {
	db := database.CreateConnection()

	db.AutoMigrate(&models.Task{})

	http.HandleFunc("/tasks", handlers.HandleTask)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
