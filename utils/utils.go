package utils

import (
	"crypto/sha1"
	"github.com/google/uuid"
)

func ToPtr[T any](arg T) *T {
	return &arg
}

func CalculateHash(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s))
	return uuid.NewSHA1(uuid.Nil, hash[:])
}
