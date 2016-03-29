// +build !linux

package gui

type DesktopNotifications struct {}

func newDesktopNotifications() *DesktopNotifications {
	return nil
}

func (dn *DesktopNotifications) show(jid, from, message string, showMessage, showFullscreen bool) error {
	return nil
}
