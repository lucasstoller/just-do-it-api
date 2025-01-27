package database

import (
	"just-do-it-api/models"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MockDB struct {
	tasks         []models.Task
	filteredTasks []models.Task
	db            *gorm.DB
}

var mockInstance *MockDB

func NewMockDB() Database {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Initialize database with Task model
	err = db.AutoMigrate(&models.Task{})
	if err != nil {
		panic("failed to migrate database")
	}

	// Create initial tasks
	tasks := []models.Task{
		{
			ID:          "1",
			Title:       "Test Task 1",
			Description: "Description 1",
			Deadline:    time.Now().Add(24 * time.Hour),
			Completed:   false,
		},
		{
			ID:          "2",
			Title:       "Test Task 2",
			Description: "Description 2",
			Deadline:    time.Now().Add(48 * time.Hour),
			Completed:   true,
		},
	}

	mockInstance = &MockDB{
		tasks: tasks,
		db:    db,
	}

	// Insert initial tasks into database
	for _, task := range tasks {
		if err := db.Create(&task).Error; err != nil {
			panic("failed to create initial tasks")
		}
	}

	return mockInstance
}

func GetMockDB() Database {
	if mockInstance == nil {
		return NewMockDB()
	}
	return mockInstance
}

func (m *MockDB) Find(dest interface{}, conds ...interface{}) *gorm.DB {
	if len(m.filteredTasks) > 0 {
		if tasks, ok := dest.(*[]models.Task); ok {
			*tasks = m.filteredTasks
			m.filteredTasks = nil
		}
		return m.db
	}
	return m.db.Find(dest, conds...)
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	if len(m.filteredTasks) > 0 {
		if task, ok := dest.(*models.Task); ok && len(m.filteredTasks) > 0 {
			*task = m.filteredTasks[0]
			m.filteredTasks = nil
			return m.db
		}
	}
	return m.db.First(dest, conds...)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	if task, ok := value.(*models.Task); ok {
		if task.ID == "" {
			task.ID = time.Now().Format("20060102150405")
		}
		m.tasks = append(m.tasks, *task)
	}
	return m.db.Create(value)
}

func (m *MockDB) Save(value interface{}) *gorm.DB {
	if task, ok := value.(*models.Task); ok {
		for i, t := range m.tasks {
			if t.ID == task.ID {
				m.tasks[i] = *task
				break
			}
		}
	}
	return m.db.Save(value)
}

func (m *MockDB) Delete(value interface{}, conds ...interface{}) *gorm.DB {
	result := m.db.Delete(value, conds...)
	if result.Error == nil && result.RowsAffected > 0 {
		if len(conds) >= 2 {
			if id, ok := conds[1].(string); ok {
				for i, task := range m.tasks {
					if task.ID == id {
						m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
						break
					}
				}
			}
		}
	}
	return result
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.filteredTasks = nil
	return m.db.Where(query, args...)
}

func (m *MockDB) AutoMigrate(dst ...interface{}) error {
	return m.db.AutoMigrate(dst...)
}
