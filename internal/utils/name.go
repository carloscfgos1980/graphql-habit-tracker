package utils

import (
	"errors"
	"strings"
	"unicode"
)

func ValidateName(name string) error {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return errors.New("name cannot be empty")
	}
	if len(name) < 2 {
		return errors.New("name must be at least 2 characters long")
	}
	if len(name) > 100 {
		return errors.New("name cannot be longer than 100 characters")
	}

	for _, char := range name {
		if unicode.IsLetter(char) {
			continue
		}
		if char == ' ' || char == '-' || char == '_' || char == '\'' || char == '.' {
			continue
		}
		return errors.New("name can only contains letters, spaces, hyphens, underscores, apostrophes, and periods")
	}

	return nil
}
