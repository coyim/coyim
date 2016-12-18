package gui

import "github.com/kardianos/osext"

// This file contains functionality for automatic installation of icons and .desktop files on platforms that support these features
// We should probably warn if we're replacing a desktop file with a different Exec path
// https://standards.freedesktop.org/desktop-entry-spec/latest/ar01s05.html

func ensureInstalled() {
	ensureIconsInstalled()
	ensureDesktopFileInstalled()
}

func ensureIconsInstalled() {
	// if a local XDG_DATA_DIRS can be found
	//    install in icons/hicolor/48x48/apps an icon with a hash of the content as part of the name
	//    install other sizes if we have
}

func ensureDesktopFileInstalled() {
	// im.coy.CoyIM.desktop
	// im.coy.CoyIM-debug.desktop
}

func generateDesktopFile(debug bool) string {
	name := "CoyIM"
	path, _ := osext.Executable()
	icon := "... TODO find icon name"
	revision := coyimVersion

	if debug {
		name = name + "-debug"
		path = path + " -debug"
	}

	return "" +
		"[Desktop Entry]\n" +
		"Type=Application\n" +
		"Encoding=UTF-8\n" +
		"Name=" + name + "\n" +
		"Comment=Secure Instant Messenger\n" +
		"Exec=" + path + "\n" +
		"Icon=" + icon + "\n" +
		"Terminal=false\n" +
		"Categories=Network;\n" +
		"# CoyIMVersion=" + revision + "\n"
}
