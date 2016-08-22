package gui

import (
	"strings"

	"github.com/twstrike/coyim/gui/settings"
	"github.com/twstrike/coyim/ui"
)

// This file contains the generic functionality for desktop notifications.
// It depends on the desktopNotifications struct having at least these three
// fields:
//     notificationStyle   string
//     notificationUrgent  bool
//     notificationExpires bool

// All platform dependent desktop notifications files should also
// contain a notificationFeaturesSupported constant that specifies what features are available

type notificationFeature int

const (
	notificationStyles notificationFeature = 2 << iota
	notificationUrgency
	notificationExpiry
)

func createDesktopNotifications() *desktopNotifications {
	dn := new(desktopNotifications)

	dn.notificationStyle = "with-author-but-no-content"
	dn.notificationUrgent = true
	dn.notificationExpires = false

	return dn
}

func (dn *desktopNotifications) updateWith(s *settings.Settings) {
	dn.notificationStyle = s.GetNotificationStyle()
	dn.notificationUrgent = s.GetNotificationUrgency()
	dn.notificationExpires = s.GetNotificationExpires()
}

func (dn *desktopNotifications) format(from, message string, withHTML bool) (summary, body string) {
	switch dn.notificationStyle {
	case "only-presence-of-new-information":
		return "New message!", ""
	case "with-author-but-no-content":
		if withHTML {
			return "New message!", "From: <b>" + from + "</b>"
		}
		return "New message!", "From: " + from
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
