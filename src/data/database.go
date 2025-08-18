package data

import (
	"alexchatapp/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func Initialize(dsn string) (*gorm.DB, error) {
	var err error
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
