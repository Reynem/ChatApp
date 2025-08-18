package utils

import (
	"alexchatapp/src/models"
	"errors"
	"regexp"
)

// ValidateEmail validates email address format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// ValidateUsername validates username format
func ValidateUsername(username string) error {
	if len(username) < 3 {
		return errors.New("username must contain at least 3 characters")
	}
	if len(username) > 50 {
		return errors.New("username must not exceed 50 characters")
	}

	// Check that username contains only allowed characters
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !usernameRegex.MatchString(username) {
		return errors.New("username can only contain letters, numbers, hyphens and underscores")
	}

	return nil
}

// ValidatePassword validates password strength
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must contain at least 8 characters")
	}
	if len(password) > 128 {
		return errors.New("password must not exceed 128 characters")
	}

	// Check for at least one digit
	hasDigit := regexp.MustCompile(`\d`).MatchString(password)
	if !hasDigit {
		return errors.New("password must contain at least one digit")
	}

	// Check for at least one letter
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	if !hasLetter {
		return errors.New("password must contain at least one letter")
	}

	return nil
}

// SanitizeUser removes sensitive information from user structure
func SanitizeUser(user *models.User) *models.User {
	if user == nil {
		return nil
	}

	sanitized := *user
	sanitized.Password = "" // Never return password
	return &sanitized
}
