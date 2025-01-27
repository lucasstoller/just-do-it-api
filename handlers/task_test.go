package handlers

import (
	"bytes"
	"encoding/json"
	"just-do-it-api/database"
	"just-do-it-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupTest(t *testing.T) {
	mockDB := database.NewMockDB()
	database.SetTestDB(mockDB)
}

func TestGetTasks(t *testing.T) {
	setupTest(t)
	// Create a request to pass to our handler
	req, err := http.NewRequest("GET", "/v1/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTasks)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	var response TaskResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	// Verify response
	if len(response.Tasks) == 0 {
		t.Error("expected tasks in response, got empty array")
	}
}

func TestCreateTask(t *testing.T) {
	setupTest(t)
	// Create test cases
	tests := []struct {
		name         string
		task         models.Task
		expectedCode int
		expectError  bool
	}{
		{
			name: "Valid Task",
			task: models.Task{
				Title:       "New Task",
				Description: "New Description",
				Deadline:    time.Now().Add(24 * time.Hour),
			},
			expectedCode: http.StatusCreated,
			expectError:  false,
		},
		{
			name: "Missing Title",
			task: models.Task{
				Description: "New Description",
				Deadline:    time.Now().Add(24 * time.Hour),
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
		{
			name: "Missing Deadline",
			task: models.Task{
				Title:       "New Task",
				Description: "New Description",
			},
			expectedCode: http.StatusBadRequest,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert task to JSON
			taskJSON, err := json.Marshal(tt.task)
			if err != nil {
				t.Fatal(err)
			}

			// Create request
			req, err := http.NewRequest("POST", "/v1/tasks", bytes.NewBuffer(taskJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(CreateTask)

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			// Check response
			if !tt.expectError {
				var response models.Task
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}
				if response.Title != tt.task.Title {
					t.Errorf("handler returned unexpected title: got %v want %v",
						response.Title, tt.task.Title)
				}
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	setupTest(t)
	// Create test cases
	tests := []struct {
		name         string
		taskID       string
		updates      models.Task
		expectedCode int
	}{
		{
			name:   "Valid Update",
			taskID: "1",
			updates: models.Task{
				Title:       "Updated Task",
				Description: "Updated Description",
				Deadline:    time.Now().Add(24 * time.Hour),
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "Task Not Found",
			taskID: "999",
			updates: models.Task{
				Title:       "Updated Task",
				Description: "Updated Description",
				Deadline:    time.Now().Add(24 * time.Hour),
			},
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert updates to JSON
			updatesJSON, err := json.Marshal(tt.updates)
			if err != nil {
				t.Fatal(err)
			}

			// Create request
			req, err := http.NewRequest("PUT", "/v1/tasks/"+tt.taskID, bytes.NewBuffer(updatesJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(UpdateTask)

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}
		})
	}
}

func TestToggleTask(t *testing.T) {
	setupTest(t)
	// Create test cases
	tests := []struct {
		name         string
		taskID       string
		expectedCode int
	}{
		{
			name:         "Valid Toggle",
			taskID:       "1",
			expectedCode: http.StatusOK,
		},
		{
			name:         "Task Not Found",
			taskID:       "999",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("PATCH", "/v1/tasks/"+tt.taskID+"/toggle", nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(ToggleTask)

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}

			// For successful toggle, check the response
			if tt.expectedCode == http.StatusOK {
				var response map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&response)
				if err != nil {
					t.Fatal(err)
				}

				if _, ok := response["completed"]; !ok {
					t.Error("response missing 'completed' field")
				}
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	setupTest(t)
	// Create test cases
	tests := []struct {
		name         string
		taskID       string
		expectedCode int
	}{
		{
			name:         "Valid Delete",
			taskID:       "1",
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "Task Not Found",
			taskID:       "999",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("DELETE", "/v1/tasks/"+tt.taskID, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Create response recorder
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(DeleteTask)

			// Call handler
			handler.ServeHTTP(rr, req)

			// Check status code
			if status := rr.Code; status != tt.expectedCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedCode)
			}
		})
	}
}

func TestGetTodayTasks(t *testing.T) {
	setupTest(t)
	// Create request
	req, err := http.NewRequest("GET", "/v1/tasks/today", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetTodayTasks)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response format
	var response TaskResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetBacklogTasks(t *testing.T) {
	setupTest(t)
	// Create request
	req, err := http.NewRequest("GET", "/v1/tasks/backlog", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create response recorder
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBacklogTasks)

	// Call handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check response format
	var response TaskResponse
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
}
