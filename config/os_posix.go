// +build !windows

package config

// IsWindows returns true if this is running under windows
func IsWindows() bool {
	return false
}

// SystemConfigDir returns the application data directory, valid on both windows and posix systems
func SystemConfigDir() string {
	return XdgConfigDir()
}
