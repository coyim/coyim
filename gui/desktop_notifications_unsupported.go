// +build !linux

package gui

type desktopNotifications struct{}

func newDesktopNotifications() *desktopNotifications {
	return nil
}

func (dn *desktopNotifications) show(jid, from, message string, showMessage, showFullscreen bool) error {
	return nil
}
