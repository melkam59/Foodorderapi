package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func Hashpassword(password string) (string, error) {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		return "", fmt.Errorf("could not hash password %w", err)
	}
	return string(hashedPassword), nil

}

/*
func Verifypassword(hashedPassword string, candidatepassword string) error {

	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatepassword))
}
*/

func VerifyPassword(hashedPassword string, candidatepassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(candidatepassword))
	if err != nil {
		return false
	}
	return true
}


