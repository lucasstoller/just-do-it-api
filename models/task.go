package models

import (
	"github.com/go-playground/validator"
	"gorm.io/gorm"
)

type TaskStatus string

const (
	Pending   TaskStatus = "pending"
	Completed TaskStatus = "completed"
)

var validate = validator.New()

type Task struct {
	gorm.Model
	Title       string     `gorm:"type:varchar(255);not null" validate:"required"`
	Description string     `gorm:"type:text"`
	Status      TaskStatus `gorm:"type:varchar(20);default:'pending'" validate:"oneof=pending completed"`
}

func (t *Task) Validate() error {
	return validate.Struct(t)
}
