// +build !windows

package config

// IsWindows returns true if this is running under windows
func IsWindows() bool {
	return false
}

// SystemConfigDir returns the application data directory, valid on both windows and posix systems
func SystemConfigDir() string {
	//TODO: Why not use g_get_user_config_dir()?
	// https://developer.gnome.org/glib/unstable/glib-Miscellaneous-Utility-Functions.html#g-get-user-config-dir
	return XdgConfigHome()
}
