package models

import (
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

var validate = validator.New()

type Task struct {
	ID          string         `gorm:"primarykey;type:varchar(255)" json:"id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title" validate:"required"`
	Description string         `gorm:"type:text" json:"description"`
	Deadline    time.Time      `gorm:"not null" json:"deadline" validate:"required"`
	Completed   bool           `gorm:"default:false" json:"completed"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = fmt.Sprintf("%d", time.Now().UnixMilli())
	}
	return nil
}

func (t *Task) Validate() error {
	return validate.Struct(t)
}
