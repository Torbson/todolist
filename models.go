package main

import (
	"gorm.io/gorm"
)

// MODELS
type Todo struct {
	gorm.Model //`json:"-"`
	//Id           int32     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	//CreationDate time.Time `json:"creation_date,omitempty"`
	//DueDate      time.Time `json:"due_date,omitempty"`
	Tasks []Task `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tasks,omitempty"`
}
type Task struct {
	gorm.Model //`json:"-"`
	//Id          int32  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	TodoID      uint
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
