package data

import (
	"gorm.io/gorm"
)

type ProfilesRepository struct {
	db *gorm.DB
}

func NewProfilesRepository(_db *gorm.DB) *ProfilesRepository {
	return &ProfilesRepository{db: _db}
}

func (r *ProfilesRepository) CreateProfile() {

}
