package pwd

import (
	"golang.org/x/crypto/bcrypt"
)

type HashService struct{}

func NewHashService() HashService {
	return HashService{}
}

func (s HashService) Hash(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), err
}

func (s HashService) Match(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
