package handlers

import (
	"encoding/json"
	"just-do-it-api/database"
	"just-do-it-api/middleware"
	"just-do-it-api/models"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"
)

type TaskResponse struct {
	Tasks []models.Task `json:"tasks"`
}

func GetTasks(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	var tasks []models.Task

	userID := middleware.GetUserID(r)
	deadlineStr := r.URL.Query().Get("deadline")

	var result *gorm.DB
	query := db.Where("user_id = ?", userID)

	if deadlineStr != "" {
		layout := "2006-01-02"
		deadline, err := time.Parse(layout, deadlineStr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(models.NewErrorResponse(
				"Invalid deadline format",
				"Deadline must be in the format YYYY-MM-DD",
			))
			return
		}

		result = query.Where("DATE(deadline) = ?", deadline.Format("2006-01-02")).Find(&tasks)
	} else {
		result = query.Find(&tasks)
	}

	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to fetch tasks",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TaskResponse{Tasks: tasks})
}

func CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	userID := middleware.GetUserID(r)

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Invalid JSON format",
		))
		return
	}
	defer r.Body.Close()

	if err := task.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			err.Error(),
		))
		return
	}

	db := database.CreateConnection()
	task.UserID = userID
	result := db.Create(&task)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to create task",
		))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Task ID is required",
		))
		return
	}

	var updates models.Task
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Invalid JSON format",
		))
		return
	}
	defer r.Body.Close()

	db := database.CreateConnection()
	var task models.Task
	userID := middleware.GetUserID(r)
	if err := db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Not found",
			"Task not found",
		))
		return
	}

	task.Title = updates.Title
	task.Description = updates.Description
	task.Deadline = updates.Deadline

	if err := task.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			err.Error(),
		))
		return
	}

	if err := db.Save(&task).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to update task",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

func DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Task ID is required",
		))
		return
	}

	db := database.CreateConnection()
	var task models.Task
	userID := middleware.GetUserID(r)
	if err := db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Not found",
			"Task not found",
		))
		return
	}

	result := db.Delete(&task)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to delete task",
		))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func ToggleTask(w http.ResponseWriter, r *http.Request) {
	taskID := strings.TrimPrefix(r.URL.Path, "/v1/tasks/")
	taskID = strings.TrimSuffix(taskID, "/toggle")
	if taskID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Invalid request",
			"Task ID is required",
		))
		return
	}

	db := database.CreateConnection()
	var task models.Task
	userID := middleware.GetUserID(r)
	if err := db.Where("id = ? AND user_id = ?", taskID, userID).First(&task).Error; err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Not found",
			"Task not found",
		))
		return
	}

	task.Completed = !task.Completed
	if err := db.Save(&task).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to toggle task",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":        task.ID,
		"completed": task.Completed,
	})
}

func GetTodayTasks(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	var tasks []models.Task

	userID := middleware.GetUserID(r)
	today := time.Now().UTC()
	startOfDay := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, time.UTC)
	endOfDay := startOfDay.Add(24 * time.Hour)

	result := db.Where("user_id = ? AND deadline BETWEEN ? AND ?", userID, startOfDay, endOfDay).Find(&tasks)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to fetch today's tasks",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TaskResponse{Tasks: tasks})
}

func GetBacklogTasks(w http.ResponseWriter, r *http.Request) {
	db := database.CreateConnection()
	var tasks []models.Task

	userID := middleware.GetUserID(r)
	now := time.Now().UTC()
	result := db.Where("user_id = ? AND deadline < ? AND completed = ?", userID, now, false).Find(&tasks)
	if result.Error != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(models.NewErrorResponse(
			"Internal server error",
			"Failed to fetch backlog tasks",
		))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TaskResponse{Tasks: tasks})
}
