package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) ([]byte, error) {
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return []byte{}, err
	}
	return hashedPass, nil
}

func CheckPasswordHash(password string, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}
