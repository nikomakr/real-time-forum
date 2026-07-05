package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)
// Despite OWASP recommends the bcrypt.DefaultCost aka 10. I have decided to boost protection by using 12. Still meets the 250ms targetn and balances protection against database leaks. 
const BcryptCost = 12 

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}
	return nil
}