package handlers

import (
	"encoding/json"
	"just-do-it-api/auth"
	"just-do-it-api/database"
	"just-do-it-api/models"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var (
	validate = validator.New()
	db       database.Database
)

func init() {
	db = database.CreateConnection()
	db.AutoMigrate(&models.User{})
}

func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Failed to parse request body",
		))
		return
	}

	if err := validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Validation error",
			err.Error(),
		))
		return
	}

	// Check if user already exists
	var existingUser models.User
	result := db.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Registration failed",
			"Email already registered",
		))
		return
	}

	// Create new user
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	if err := user.HashPassword(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Registration failed",
			"Failed to hash password",
		))
		return
	}

	if result := db.Create(&user); result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Registration failed",
			"Failed to create user",
		))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Registration failed",
			"Failed to generate token",
		))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  user,
	})
}

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Failed to parse request body",
		))
		return
	}

	if err := validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Validation error",
			err.Error(),
		))
		return
	}

	// Find user by email
	var user models.User
	result := db.Where("email = ?", req.Email).First(&user)
	if result.Error != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Login failed",
			"Invalid email or password",
		))
		return
	}

	// Check password
	if err := user.CheckPassword(req.Password); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Login failed",
			"Invalid email or password",
		))
		return
	}

	// Generate JWT token
	token, err := auth.GenerateToken(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Login failed",
			"Failed to generate token",
		))
		return
	}

	json.NewEncoder(w).Encode(models.AuthResponse{
		Token: token,
		User:  user,
	})
}
