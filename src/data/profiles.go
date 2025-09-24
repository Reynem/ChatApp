package data

import (
	"alexchatapp/src/models"
	"time"

	"gorm.io/gorm"
)

type ProfilesRepository struct {
	db *gorm.DB
}

func NewProfilesRepository(_db *gorm.DB) *ProfilesRepository {
	return &ProfilesRepository{db: _db}
}

func (r *ProfilesRepository) CreateProfile(user_id uint, profile_name string) error {
	var new_profile = models.Profile{
		User_id:      user_id,
		Profile_name: profile_name,
	}

	return r.db.Create(new_profile).Error
}

func (r *ProfilesRepository) UpdateProfile(profile *models.Profile) error {
	profile.Last_seen = time.Now()
	return r.db.Save(profile).Error
}

func (r *ProfilesRepository) DeleteProfile(user_id uint) error {
	return r.db.Delete(&models.Profile{}, user_id).Error
}

func (r *ProfilesRepository) GetProfileByID(user_id uint) (*models.Profile, error) {
	var profile models.Profile
	err := r.db.First(&profile, user_id).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *ProfilesRepository) DoesProfileExist(user_id uint) bool {
	var profile models.Profile
	err := r.db.First(&profile, user_id).Error
	return err == nil
}
