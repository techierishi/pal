package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func CalculateHash(input string) string {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)
	return hex.EncodeToString(hash)
}
