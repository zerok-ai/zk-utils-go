package crypto

import (
	"crypto/sha1"
	"github.com/google/uuid"
)

func CalculateHashNewSHA2(s string) uuid.UUID {
	hash := sha1.Sum([]byte(s))
	return uuid.NewSHA1(uuid.Nil, hash[:])
}
