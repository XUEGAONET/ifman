package main

import "crypto/sha256"

func keyHash(s string) []byte {
	hash := sha256.Sum256([]byte(s))
	t := make([]byte, 32)
	copy(t, hash[:])
	return t
}
