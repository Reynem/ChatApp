package data

import (
	"errors"
	"time"

	"alexchatapp/src/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthRepository contains methods for authentication database operations
type AuthRepository struct {
	db *gorm.DB
}

// NewAuthRepository creates a new authentication repository instance
func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db: db}
}

// CreateUser creates a new user in the database
func (r *AuthRepository) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

// GetUserByUsername finds a user by username
func (r *AuthRepository) GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByEmail finds a user by email
func (r *AuthRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByID finds a user by ID
func (r *AuthRepository) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user data
func (r *AuthRepository) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	return r.db.Save(user).Error
}

// DeleteUser deletes a user
func (r *AuthRepository) DeleteUser(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UserExists checks if a user exists by username
func (r *AuthRepository) UserExists(username string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("user_name = ?", username).Count(&count).Error
	return count > 0, err
}

// EmailExists checks if a user exists by email
func (r *AuthRepository) EmailExists(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// HashPassword hashes a password
func (r *AuthRepository) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies a password
func (r *AuthRepository) CheckPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// RegisterUser registers a new user
func (r *AuthRepository) RegisterUser(username, email, password string) (*models.User, error) {
	// Check if user exists
	exists, err := r.UserExists(username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("user with this username already exists")
	}

	// Check if email exists
	emailExists, err := r.EmailExists(email)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := r.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		UserName:  username,
		Email:     email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := r.CreateUser(user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser authenticates a user
func (r *AuthRepository) AuthenticateUser(username, password string) (*models.User, error) {
	// Find user
	user, err := r.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Check password
	if err := r.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
