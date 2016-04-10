// +build !linux

package gui

import "github.com/twstrike/coyim/gui/settings"

type desktopNotifications struct{}

func newDesktopNotifications() *desktopNotifications {
	return nil
}

func (dn *desktopNotifications) updateWith(*settings.Settings) {
}

func (dn *desktopNotifications) show(jid, from, message string) error {
	return nil
}
