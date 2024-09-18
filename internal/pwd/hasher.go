package pwd

import (
	"golang.org/x/crypto/bcrypt"
)

// HashService provides methods for hashing passwords and comparing hashed passwords.
type HashService struct{}

// NewHashService creates a new instance of HashService.
func NewHashService() HashService {
	return HashService{}
}

// Hash takes a plain text password and returns its bcrypt hash.
// It uses bcrypt's MinCost parameter for hashing, which is a low computational cost for fast hashing.
// Returns the hashed password as a string and any error encountered during hashing.
func (s HashService) Hash(password string) (string, error) {
	var passwordBytes = []byte(password)
	hashedPasswordBytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.MinCost)
	return string(hashedPasswordBytes), err
}

// Match compares a plain text password with a hashed password to check if they match.
// It returns true if the password matches the hash, and false otherwise.
// Returns true if the plain text password, when hashed, matches the provided hashed password, and false if not.
func (s HashService) Match(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
