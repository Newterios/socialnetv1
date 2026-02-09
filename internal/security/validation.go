package security

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("invalid email format")
	}
	return nil
}

func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if len(username) > 30 {
		return errors.New("username must be at most 30 characters")
	}
	if !isAlphanumeric(username) {
		return errors.New("username must contain only letters, numbers, and underscores")
	}
	return nil
}

func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	if len(password) > 100 {
		return errors.New("password must be at most 100 characters")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false

	for _, char := range password {
		if unicode.IsUpper(char) {
			hasUpper = true
		}
		if unicode.IsLower(char) {
			hasLower = true
		}
		if unicode.IsDigit(char) {
			hasDigit = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit {
		return errors.New("password must contain uppercase, lowercase, and digit")
	}

	return nil
}

func ValidateContent(content string, maxLength int) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("content is required")
	}
	if len(content) > maxLength {
		return errors.New("content exceeds maximum length")
	}
	return nil
}

func isAlphanumeric(s string) bool {
	for _, char := range s {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && char != '_' {
			return false
		}
	}
	return true
}
