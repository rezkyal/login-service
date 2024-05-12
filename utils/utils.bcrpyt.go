package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func VerifyPassword(hashedPassword, plainPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))

	// wrong password
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	if err != nil {
		return false, fmt.Errorf("[VerifyPassword] error when CompareHashAndPassword, err: %+v", err)
	}

	return true, nil
}
