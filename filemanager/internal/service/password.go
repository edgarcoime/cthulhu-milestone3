package service

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = bcrypt.DefaultCost

// HashBucketPassword hashes a plaintext bucket password for storage.
func HashBucketPassword(password string) (string, error) {
	if password == "" {
		return "", fmt.Errorf("password cannot be empty")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(hash), nil
}

// VerifyBucketPassword returns true if the plaintext password matches the hash.
func VerifyBucketPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
