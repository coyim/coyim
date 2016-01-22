package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) showConnectAccountNotification(account *account) func() {
	var notification *gtk.InfoBar

	doInUIThread(func() {
		notification = account.buildConnectionNotification()
		account.setCurrentNotification(notification, u.notificationArea)
	})

	return func() {
		doInUIThread(func() {
			account.removeCurrentNotificationIf(notification)
		})
	}
}

func (u *gtkUI) notifyConnectionFailure(account *account) {
	doInUIThread(func() {
		notification := account.buildConnectionFailureNotification()
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func buildVerifyIdentityNotification(acc *account, peer string, win *gtk.Window) *gtk.InfoBar {
	builder := builderForDefinition("VerifyIdentityNotification")

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	message := obj.(*gtk.Label)
	message.SetSelectable(true)

	text := fmt.Sprintf(i18n.Local("You have not verified the identity of %s"), peer)
	message.SetText(text)

	obj, _ = builder.GetObject("button_verify")
	button := obj.(*gtk.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
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

func (u *gtkUI) notify(title, message string) {
	builder := builderForDefinition("SimpleNotification")
	obj, _ := builder.GetObject("dialog")
	dlg := obj.(*gtk.MessageDialog)

	dlg.SetProperty("title", title)
	dlg.SetProperty("text", message)
	dlg.SetTransientFor(u.window)

	doInUIThread(func() {
		dlg.Run()
		dlg.Destroy()
	})
}
