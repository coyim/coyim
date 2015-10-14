package config

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"strings"
)

func ParseYes(input string) bool {
	switch strings.ToLower(input) {
	case "y", "yes":
		return true
	}

	return false
}

func randomString(dest []byte) error {
	src := make([]byte, len(dest))

	if _, err := io.ReadFull(rand.Reader, src); err != nil {
		return err
	}

	copy(dest, hex.EncodeToString(src))

	return nil
}
