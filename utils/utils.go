package utils

import (
	"crypto/sha1"
	"github.com/google/uuid"
)

func ToPtr[T any](arg T) *T {
	return &arg
}

func CalculateHash(s string) uuid.UUID {
	// Calculate the SHA-1 hash of the sorted JSON string
	hash := sha1.Sum([]byte(s))

	// Create a UUID from the hash
	return uuid.NewSHA1(uuid.Nil, hash[:])
}
