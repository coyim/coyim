// +build !windows

package config

import "path/filepath"

func configDir() string {
	return filepath.Join(xdgHomeDir(), "coyim")
}
