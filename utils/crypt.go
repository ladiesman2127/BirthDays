package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password *string) (*string, error) {
	hashedPassInBytes, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashedPassword := string(hashedPassInBytes)
	return &hashedPassword, nil
}
