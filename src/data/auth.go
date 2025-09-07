package data

import (
	"alexchatapp/src/models"
	"alexchatapp/src/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

// RegisterUser registers a new user
func (r *UsersRepository) RegisterUser(username, email, password string) (*models.User, error) {
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
	hashedPassword, err := utils.HashPassword(password)
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

// AuthenticateUser authenticates a user (Login)
func (r *UsersRepository) AuthenticateUser(username, password string) (*models.User, error) {
	// Find user
	user, err := r.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid username or password")
		}
		return nil, err
	}

	// Check password
	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
