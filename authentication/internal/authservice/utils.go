package authservice

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func encryptPassword(password string) (string, error) {
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("error generating password %w", err)
	}
	return string(encrypted), nil
}

func compareHashAndPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
