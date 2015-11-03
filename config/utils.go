package config

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ParseYes returns true if the string is any combination of yes
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

func xdgOr(env, or string) string {
	x := os.Getenv(env)
	if x == "" {
		x = filepath.Join(os.Getenv("HOME"), or)
	}
	return x
}

// XdgConfigDir returns the standardized XDG Configuration directory
func XdgConfigDir() string {
	return xdgOr("XDG_CONFIG_HOME", ".config")
}

// XdgCacheDir returns the standardized XDG Cache directory
func XdgCacheDir() string {
	return xdgOr("XDG_CACHE_HOME", ".cache")
}

// XdgDataDir returns the standardized XDG Data directory
func XdgDataDir() string {
	return xdgOr("XDG_DATA_HOME", ".local/share")
}
