// +build linux

package gui

import (
	"fmt"
	"log"
	"strings"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/TheCreeper/go-notify"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/godbus/dbus"
	"github.com/twstrike/coyim/gui/settings"

	"github.com/twstrike/coyim/ui"
)

type desktopNotifications struct {
	notifications       map[string]uint32
	notification        notify.Notification
	notificationStyle   string
	notificationUrgent  bool
	notificationExpires bool
}

func (dn *desktopNotifications) updateWith(s *settings.Settings) {
	dn.notificationStyle = s.GetNotificationStyle()
	dn.notificationUrgent = s.GetNotificationUrgency()
	dn.notificationExpires = s.GetNotificationExpires()
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

	dn := new(desktopNotifications)
	dn.notifications = make(map[string]uint32)
	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	return dn
}

func (dn *desktopNotifications) format(from, message string) (summary, body string) {
	switch dn.notificationStyle {
	case "only-presence-of-new-information":
		return "New message!", ""
	case "with-author-but-no-content":
		return "New message!", "From: <b>" + from + "</b>"
	case "with-content":
		smsg := strings.Split(message, "\n")[0]
		smsg = ui.EscapeAllHTMLTags(smsg)
		if len(smsg) > 254 {
			smsg = smsg[0:253]
			stok := strings.Split(smsg, " ")
			if len(stok) > 1 {
				smsg = strings.Join(stok[0:len(stok)-2], " ")
			}
			smsg = smsg + "..."
		}
		return "From: " + from, smsg
	}
	return "", ""
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
	notification.Summary, notification.Body = dn.format(from, message)

	nid, err := notification.Show()
	if err != nil {
		return fmt.Errorf("Error showing notification: %v", err)
	}
	dn.notifications[jid] = nid
	dn.notification = notification
	return nil
}
