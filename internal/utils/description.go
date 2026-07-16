package utils

import (
	"errors"
	"strings"
)

func ValidateDescription(description string) error {
	description = strings.TrimSpace(description)

	if len(description) == 0 {
		return errors.New("description cannot be empty")
	}

	if len(description) < 10 {
		return errors.New("description must be at least 10 characters long")
	}

	if len(description) > 250 {
		return errors.New("description must not exceed 250 characters")
	}

	return nil
}
