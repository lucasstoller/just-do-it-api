package handlers

import (
	"bytes"
	"encoding/json"
	"just-do-it-api/database"
	"just-do-it-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type AuthMockDB struct {
	db *gorm.DB
}

func NewAuthMockDB() database.Database {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Initialize database with User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic("failed to migrate database")
	}

	return &AuthMockDB{db: db}
}

func (m *AuthMockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	return m.db.Find(dest, conds...)
}

func (m *AuthMockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	return m.db.First(dest, conds...)
}

func (m *AuthMockDB) Create(value interface{}) *gorm.DB {
	return m.db.Create(value)
}

func (m *AuthMockDB) Save(value interface{}) *gorm.DB {
	return m.db.Save(value)
}

func (m *AuthMockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	return m.db.Delete(value, conds...)
}

func (m *AuthMockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	return m.db.Where(query, args...)
}

func (m *AuthMockDB) AutoMigrate(dst ...interface{}) error {
	return m.db.AutoMigrate(dst...)
}

func TestRegister(t *testing.T) {
	mockDB := NewAuthMockDB()
	db = mockDB

	tests := []struct {
		name           string
		payload        models.RegisterRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid Registration",
			payload: models.RegisterRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "Invalid Email",
			payload: models.RegisterRequest{
				Email:    "invalid-email",
				Password: "password123",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name: "Short Password",
			payload: models.RegisterRequest{
				Email:    "test@example.com",
				Password: "short",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader(payloadBytes))
			w := httptest.NewRecorder()

			Register(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			if tt.expectedError {
				if _, ok := response["error"]; !ok {
					t.Error("expected error response, got success")
				}
			} else {
				if token, ok := response["token"].(string); !ok || token == "" {
					t.Error("expected token in response")
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	mockDB := NewAuthMockDB()
	db = mockDB

	// Create a test user
	testUser := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}
	testUser.HashPassword()
	if err := mockDB.(*AuthMockDB).db.Create(testUser).Error; err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	tests := []struct {
		name           string
		payload        models.LoginRequest
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "Valid Login",
			payload: models.LoginRequest{
				Email:    "test@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name: "Invalid Email",
			payload: models.LoginRequest{
				Email:    "nonexistent@example.com",
				Password: "password123",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
		{
			name: "Wrong Password",
			payload: models.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			expectedStatus: http.StatusUnauthorized,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			payloadBytes, _ := json.Marshal(tt.payload)
			req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader(payloadBytes))
			w := httptest.NewRecorder()

			Login(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response map[string]interface{}
			json.NewDecoder(w.Body).Decode(&response)

			if tt.expectedError {
				if _, ok := response["error"]; !ok {
					t.Error("expected error response, got success")
				}
			} else {
				tokenInterface, ok := response["token"]
				if !ok {
					t.Error("expected token in response")
					return
				}
				if tokenStr, ok := tokenInterface.(string); !ok || tokenStr == "" {
					t.Error("expected token to be non-empty string")
				}
			}
		})
	}
}
