// +build linux

package gui

import (
	"fmt"
	"log"
	"strings"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/TheCreeper/go-notify"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/godbus/dbus"

	"github.com/twstrike/coyim/ui"
)

type desktopNotifications struct {
	notifications map[string]uint32
	notification notify.Notification
}

func newDesktopNotifications() *desktopNotifications {
	if _, err := dbus.SessionBus(); err != nil {
		log.Printf("Error enabling dbus based notifications! %+v\n", err)
		return nil
	}

	dn := new(desktopNotifications)
	dn.notifications = make(map[string]uint32)
	return dn
}

func (dn *desktopNotifications) show(jid, from, message string, showMessage, showFullscreen bool) error {
	hints := make(map[string]interface{})
	//hints[notify.HintResident] = true
	hints[notify.HintTransient] = false
	hints[notify.HintActionIcons] = "coyim"
	hints[notify.HintDesktopEntry] = "coyim.desktop"
	hints[notify.HintCategory] = notify.ClassImReceived
	if showFullscreen {
		hints[notify.HintUrgency] = notify.UrgencyCritical
	}
	notification := notify.Notification{
		AppName:    "CoyIM",
		AppIcon:    "coyim",
		Timeout:    notify.ExpiresNever,
		Hints:      hints,
		ReplacesID: dn.notifications[jid],
	}

	from = ui.EscapeAllHTMLTags(string(ui.StripSomeHTML([]byte(from))))
	if message == "" || showMessage == false {
		notification.Summary = "New message!"
		notification.Body = "From: <b>" + from + "</b>"
	} else {
		notification.Summary = "From: " + from
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
		notification.Body = smsg
	}

	nid, err := notification.Show()
	if err != nil {
		return fmt.Errorf("Error showing notification: %v", err)
	}
	dn.notifications[jid] = nid
	dn.notification = notification
	return nil
}
