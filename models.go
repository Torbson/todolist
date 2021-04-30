package main

import (
	"time"

	"gorm.io/gorm"
)

// MODELS
type Todo struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name" validate:"min=1,max=128"`
	Description string         `json:"description,omitempty"`
	Tasks       []Task         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tasks,omitempty"`
}
type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name" validate:"min=1,max=128"`
	Description string         `json:"description,omitempty"`
	TodoID      uint           `json:"-"`
}
