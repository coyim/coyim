package gui

import (
	"fmt"

	"github.com/coyim/coyim/ui"
	"github.com/coyim/gosx-notifier"
)

var hasBundle = false

const notificationFeaturesSupported = notificationStyles

func init() {
	hasBundle = !fileNotFound("/Applications/CoyIM.app/Contents/Info.plist")
}

type desktopNotifications struct {
	notificationStyle   string
	notificationUrgent  bool
	notificationExpires bool
}

func newDesktopNotifications() *desktopNotifications {
	return createDesktopNotifications()
}

func (dn *desktopNotifications) show(jid, from, message string) error {
	if dn.notificationStyle == "off" {
		return nil
	}

	from = ui.EscapeAllHTMLTags(string(ui.StripSomeHTML([]byte(from))))
	summary, body := dn.format(from, message, false)
	note := gosxnotifier.NewNotification(body)
	note.Title = summary

	note.Group = fmt.Sprintf("im.coy.coyim.%s", from)

	if hasBundle {
		note.Sender = "im.coy.coyim"
	} else {
		note.ContentImage = coyimIcon.getPath()
	}

	err := note.Push()
	if err != nil {
		return fmt.Errorf("Error showing notification: %v", err)
	}

	return nil
}
