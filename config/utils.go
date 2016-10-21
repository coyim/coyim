package config

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
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

const fingerprintDefaultGrouping = 8

// FormatFingerprint returns a formatted string of the fingerprint
func FormatFingerprint(fpr []byte) string {
	str := fmt.Sprintf("%X", fpr)
	result := ""

	sep := ""
	for len(str) > 0 {
		result = result + sep + str[0:fingerprintDefaultGrouping]
		sep = " "
		str = str[fingerprintDefaultGrouping:]
	}

	return result
}

func randomString(dest []byte) error {
	src := make([]byte, len(dest))

	if _, err := io.ReadFull(rand.Reader, src); err != nil {
		return err
	}

	copy(dest, hex.EncodeToString(src))

	return nil
}

// WithHome returns the given relative file/dir with the $HOME prepended
func WithHome(file string) string {
	return filepath.Join(os.Getenv("HOME"), file)
}

func xdgOrWithHome(env, or string) string {
	x := os.Getenv(env)
	if x == "" {
		x = WithHome(or)
	}
	return x
}

// FindFile will check each path and if that file exists return the file name and true
func FindFile(places []string) (string, bool) {
	for _, p := range places {
		if fileExists(p) {
			return p, true
		}
	}
	return "", false
}

// XdgConfigHome returns the standardized XDG Configuration directory
func XdgConfigHome() string {
	return xdgOrWithHome("XDG_CONFIG_HOME", ".config")
}

// XdgCacheDir returns the standardized XDG Cache directory
func XdgCacheDir() string {
	return xdgOrWithHome("XDG_CACHE_HOME", ".cache")
}

// XdgDataHome returns the standardized XDG Data directory
func XdgDataHome() string {
	return xdgOrWithHome("XDG_DATA_HOME", ".local/share")
}

// XdgDataDirs returns the standardized XDG Data directory
func XdgDataDirs() []string {
	x := os.Getenv("XDG_DATA_DIRS")
	return strings.Split(x, ":")
}
