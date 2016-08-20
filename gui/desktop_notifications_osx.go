// +build darwin

package gui

import "github.com/twstrike/coyim/gui/settings"
import "github.com/twstrike/gosx-notifier"

type desktopNotifications struct {
	notificationStyle   string
	notificationUrgent  bool
	notificationExpires bool
}

func (dn *desktopNotifications) updateWith(s *settings.Settings) {
	dn.notificationStyle = s.GetNotificationStyle()
	dn.notificationUrgent = s.GetNotificationUrgency()
	dn.notificationExpires = s.GetNotificationExpires()
}

func newDesktopNotifications() *desktopNotifications {
	dn := new(desktopNotifications)
	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	return dn
}

func (dn *desktopNotifications) show(jid, from, message string) error {
	if dn.notificationStyle == "off" {
		return nil
	}
	gosxnotifier.NewNotification("Check your Apple Stock!")
	return nil
}
