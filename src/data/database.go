package data

import (
	"alexchatapp/src/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// ChatRepository contains methods for database operations
type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) GetDB() *gorm.DB {
	return r.db
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
