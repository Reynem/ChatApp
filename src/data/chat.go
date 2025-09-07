package data

import (
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
