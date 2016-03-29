// +build linux

package gui

import (
	"fmt"
	"log"
	"strings"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/TheCreeper/go-notify"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/godbus/dbus"
)

type DesktopNotifications struct {
	notifications map[string]uint32
}

func newDesktopNotifications() *DesktopNotifications {
	if _, err := dbus.SessionBus(); err != nil {
		log.Printf("Error enabling dbus based notifications! %+v\n", err)
		return nil
	}

	dn := new(DesktopNotifications)
	dn.notifications = make(map[string]uint32)
	return dn
}

func (dn *DesktopNotifications) show(jid, from, message string, showMessage, showFullscreen bool) error {
	hints := make(map[string]interface{})
	//hints[notify.HintResident] = true
	hints[notify.HintTransient] = false
	hints[notify.HintActionIcons] = "coyim"
	if showFullscreen {
		hints[notify.HintUrgency] = notify.UrgencyCritical
	}
	notification := notify.Notification{
		AppName: "CoyIM",
		AppIcon: "coyim",
		Timeout: notify.ExpiresNever,
		Hints:   hints,
	}
	if message == "" || showMessage == false {
		notification.Summary = "New message!"
		notification.Body = "From: <b>" + from + "</b>"
	} else {
		notification.Summary = "From: " + from
		smsg := strings.Split(message, "\n")[0]
		if len(smsg) > 254 {
			smsg = smsg[0:253]
			stok := strings.Split(smsg, " ")
			if len(stok) > 1 {
				smsg = strings.Join(stok[0:len(stok)-2], " ")
			}
			smsg = smsg + "…"
		}
		notification.Body = smsg
	}

	nid, err := notification.Show()
	if err != nil {
		return fmt.Errorf("Error showing notification: %v", err)
	}
	dn.notifications[jid] = nid
	return nil
}
