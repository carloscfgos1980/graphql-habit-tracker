package utils

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateEmail(email string) error {
	email = strings.TrimSpace(email)

	if len(email) == 0 {
		return errors.New("email cannot be empty")
	}

	if len(email) > 254 {
		return errors.New("email address is too long")
	}

	if !emailRegex.MatchString(email) {
		return errors.New("email address is not valid")
	}

	return nil

}
