// +build darwin

package gui

import (
	"log"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gosx-notifier"
	"github.com/twstrike/coyim/gui/settings"
)

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
	note := gosxnotifier.NewNotification("Check your Apple Stock!")

	// //Optionally, set a title
	// note.Title = "It's money making time 💰"

	// //Optionally, set a subtitle
	// note.Subtitle = "My subtitle"

	// //Optionally, set a sound from a predefined set.
	// note.Sound = gosxnotifier.Basso

	// //Optionally, set a group which ensures only one notification is ever shown replacing previous notification of same group id.
	// note.Group = "com.unique.yourapp.identifier"

	// //Optionally, set a sender (Notification will now use the Safari icon)
	// note.Sender = "com.apple.Safari"

	// //Optionally, specifiy a url or bundleid to open should the notification be
	// //clicked.
	// note.Link = "http://www.yahoo.com" //or BundleID like: com.apple.Terminal

	// //Optionally, an app icon (10.9+ ONLY)
	// note.AppIcon = "gopher.png"

	// //Optionally, a content image (10.9+ ONLY)
	// note.ContentImage = "gopher.png"

	//Then, push the notification
	err := note.Push()

	//If necessary, check error
	if err != nil {
		log.Println("Uh oh!")
	}

	return nil
}
