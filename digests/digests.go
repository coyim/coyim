package digests

import (
	"crypto/sha1"
	"crypto/sha256"

	"golang.org/x/crypto/sha3"
)

// Sha1 will generate a SHA1 digest
func Sha1(i []byte) []byte {
	h := sha1.New()
	h.Write(i)
	return h.Sum(nil)
}

// Sha256 will generate a SHA256 digest
func Sha256(i []byte) []byte {
	h := sha256.New()
	h.Write(i)
	return h.Sum(nil)
}

// Sha3_256 will generate a SHA3-256 digest
func Sha3_256(i []byte) []byte {
	h := sha3.New256()
	h.Write(i)
	return h.Sum(nil)
}
