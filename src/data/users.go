package data

import (
	"time"

	"alexchatapp/src/models"

	"gorm.io/gorm"
)

// AuthRepository contains methods for authentication database operations
type UsersRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new authentication repository instance
func NewUsersRepository(db *gorm.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *UsersRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByUsername finds a user by username
func (r *UsersRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail finds a user by email
func (r *UsersRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID finds a user by ID
func (r *UsersRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user data
func (r *UsersRepository) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// DeleteUser deletes a user
func (r *UsersRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UserExists checks if a user exists by username
func (r *UsersRepository) UserExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("user_name = ?", username).Count(&count).Error
	return count > 0, err
}

// EmailExists checks if a user exists by email
func (r *UsersRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
