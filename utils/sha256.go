package utils

import "crypto/sha256"

func SHA256(s string) [32]byte {
	return sha256.Sum256([]byte(s))
}
