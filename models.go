package main

import (
	"time"

	"gorm.io/gorm"
)

// MODELS
type Todo struct {
	//gorm.Model //`json:"-"`
	//Id           int32       `json:"id"`
	ID          uint           `gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	//CreationDate time.Time `json:"creation_date,omitempty"`
	//DueDate      time.Time `json:"due_date,omitempty"`
	Tasks []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tasks,omitempty"`
}
type Task struct {
	//gorm.Model //`json:"-"`
	//Id          int32  `json:"id"`
	ID          uint           `gorm:"primaryKey"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	TodoID      uint           `json:"-"`
	//Duration string `json:"duration,omitempty"`
	//Items []Item `json:"items,omitempty"`
}

/*type Item struct {
	gorm.Model
	Name string `json:"name,omitempty"`
	Vendor string `json:"vendor"`
	Article string `json:"article"`
	Amount int32 `json:"amount,omitempty"`
	TaskID      uint
}*/
