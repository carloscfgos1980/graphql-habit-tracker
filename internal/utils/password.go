package utils

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

func ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	var (
		hasUpper  bool
		hasLower  bool
		hasNumber bool
		hasSymbol bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true

		case unicode.IsLower(char):
			hasLower = true

		case unicode.IsDigit(char):
			hasNumber = true

		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSymbol {
		return errors.New("password must contain at least one symbol")
	}

	return nil
}

func HashPassword(password string) (string, error) {
	passwordBytes := []byte(password)

	// 10 -> 100ms
	// 12 -> 400ms
	// 14 -> 1.5s
	hashedBytes, err := bcrypt.GenerateFromPassword(passwordBytes, 12)

	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil

}

func ComparePassword(hashedPassword string, plainPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
}
