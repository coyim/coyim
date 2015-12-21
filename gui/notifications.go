package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) showConnectAccountNotification(account *account) func() {
	notification := account.buildConnectionNotification()

	glib.IdleAdd(func() {
		account.setCurrentNotification(notification, u.notificationArea)
	})

	return func() {
		glib.IdleAdd(func() {
			account.removeCurrentNotificationIf(notification)
		})
	}
}

func (u *gtkUI) notifyConnectionFailure(account *account) {
	notification := account.buildConnectionFailureNotification()

	glib.IdleAdd(func() {
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func buildVerifyIdentityNotification(acc *account, peer string, win *gtk.Window) *gtk.InfoBar {
	builder := builderForDefinition("VerifyIdentityNotification")

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	message := obj.(*gtk.Label)

	text := fmt.Sprintf(i18n.Local("You have not verified the identity of %s"), peer)
	message.SetText(text)

	obj, _ = builder.GetObject("button_verify")
	button := obj.(*gtk.Button)
	button.Connect("clicked", func() {
		glib.IdleAdd(func() {
			resp := verifyFingerprintDialog(acc, peer, win)
			if resp == gtk.RESPONSE_YES {
				infoBar.Hide()
				infoBar.Destroy()
			}
		})
	})

	infoBar.ShowAll()

	return infoBar
}
