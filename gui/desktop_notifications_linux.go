package gui

import (
	"fmt"
	"log"
	"time"

	"github.com/TheCreeper/go-notify"
	"github.com/godbus/dbus"

	"github.com/coyim/coyim/ui"
)

const notificationFeaturesSupported = notificationStyles | notificationUrgency | notificationExpiry

type desktopNotifications struct {
	notifications       map[string]uint32
	notification        notify.Notification
	notificationStyle   string
	notificationUrgent  bool
	notificationExpires bool
	supported           bool
}

func (dn *desktopNotifications) hints() map[string]interface{} {
	if !dn.supported {
		return nil
	}

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
		return createDesktopNotifications()
	}

	dn := createDesktopNotifications()
	dn.supported = true
	dn.notifications = make(map[string]uint32)
	return dn
}

const defaultExpirationMs = 5000

func (dn *desktopNotifications) expiration() int32 {
	if dn.notificationExpires {
		return defaultExpirationMs
	}
	return notify.ExpiresNever
}

func (dn *desktopNotifications) show(jid, from, message string) error {
	if dn.notificationStyle == "off" || !dn.supported {
		return nil
	}

	notification := notify.Notification{
		AppName:    "CoyIM",
		AppIcon:    "coyim",
		Timeout:    dn.expiration(),
		Hints:      dn.hints(),
		ReplacesID: dn.notifications[jid],
	}

//	from = ui.EscapeAllHTMLTags(string(ui.StripSomeHTML([]byte(from))))
	from = string(ui.StripSomeHTML([]byte(from)))
	notification.Summary, notification.Body = dn.format(from, message, true)

	nid, err := notification.Show()
	if err != nil {
		return fmt.Errorf("Error showing notification: %v", err)
	}

	if dn.notificationExpires {
		go expireNotification(nid, defaultExpirationMs)
	}

	dn.notifications[jid] = nid
	dn.notification = notification
	return nil
}

func expireNotification(id uint32, expiry int) {
	time.Sleep(time.Duration(expiry) * time.Millisecond)
	notify.CloseNotification(id)
}
