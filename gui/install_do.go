//go:build !no_resource_install
// +build !no_resource_install

package gui

func ensureInstalled() {
	ensureIconsInstalled()
	ensureDesktopFileInstalled()
}
