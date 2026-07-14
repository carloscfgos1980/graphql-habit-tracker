package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"userId"`

	// ID,
	// Subject
	// Issuer
	// NotBefore
	// ExpiresAt
	// IssuedAt

	jwt.RegisteredClaims
}

func GenerateJWT(userID string, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", nil
	}

	return tokenString, nil
}

func ValidateJWT(tokenString string, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, errors.New("unexpected signing method - possible attack")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(*Claims)

	if !ok {
		return "", errors.New("could not extract claims from token")
	}

	return claims.UserID, nil

}
