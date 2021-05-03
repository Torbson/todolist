package main

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func connect_db(dbname string) {
	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", POSTGRES_HOST, POSTGRES_USER, POSTGRES_PASSWORD, dbname, POSTGRES_PORT)
	// Connect PostgreSQL
	log.Print("connect to Postgres ...")
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("failed to connect database")
	}
	log.Print("successfully connected to Postgres")
	// Auto migrate in case of schema / model changes
	db.AutoMigrate(&Todo{}, &Task{})
}

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
