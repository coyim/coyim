package gui

import (
	"fmt"
	"log"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/TheCreeper/go-notify"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/godbus/dbus"

	"github.com/twstrike/coyim/ui"
)

type desktopNotifications struct {
	notifications       map[string]uint32
	notification        notify.Notification
	notificationStyle   string
	notificationUrgent  bool
	notificationExpires bool
}

func (dn *desktopNotifications) hints() map[string]interface{} {
	hints := make(map[string]interface{})

	hints[notify.HintTransient] = false
	hints[notify.HintActionIcons] = "coyim"
	hints[notify.HintDesktopEntry] = "coyim.desktop"
	hints[notify.HintCategory] = notify.ClassImReceived

	if dn.notificationUrgent {
		hints[notify.HintUrgency] = notify.UrgencyCritical
	} else {
		hints[notify.HintUrgency] = notify.UrgencyNormal
	}
	return hints
}

func newDesktopNotifications() *desktopNotifications {
	if _, err := dbus.SessionBus(); err != nil {
		log.Printf("Error enabling dbus based notifications! %+v\n", err)
		return nil
	}

	dn := createDesktopNotifications()
	dn.notifications = make(map[string]uint32)
	return dn
}

func (dn *desktopNotifications) expiration() int32 {
	if dn.notificationExpires {
		return notify.ExpiresDefault
	}
	return notify.ExpiresNever
}

func (dn *desktopNotifications) show(jid, from, message string) error {
	if dn.notificationStyle == "off" {
		return nil
	}

	notification := notify.Notification{
		AppName:    "CoyIM",
		AppIcon:    "coyim",
		Timeout:    dn.expiration(),
		Hints:      dn.hints(),
		ReplacesID: dn.notifications[jid],
	}

	from = ui.EscapeAllHTMLTags(string(ui.StripSomeHTML([]byte(from))))
	notification.Summary, notification.Body = dn.format(from, message, true)

	nid, err := notification.Show()
	if err != nil {
		return fmt.Errorf("Error showing notification: %v", err)
	}
	dn.notifications[jid] = nid
	dn.notification = notification
	return nil
}
