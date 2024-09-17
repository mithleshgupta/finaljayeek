package security

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

const defaultCost = 12

// Hash hashes the given password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), defaultCost)
}

// VerifyPassword verifies that the given password matches the hashed password
func VerifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// IsHashed checks if the given password is already hashed
func IsHashed(password string) bool {
	// Define the regular expression pattern for bcrypt hash
	bcryptPattern := `^\$2a\$([0-9]{2})\$[./0-9A-Za-z]{53}$`

	// Create a regular expression object
	regexpHash := regexp.MustCompile(bcryptPattern)

	// Check if the password matches the bcrypt hash pattern
	return regexpHash.MatchString(password)
}
