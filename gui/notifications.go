package gui

import (
	"fmt"

	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) showConnectAccountNotification(account *account) func() {
	var notification gtki.InfoBar

	doInUIThread(func() {
		notification = account.buildConnectionNotification(u)
		account.setCurrentNotification(notification, u.notificationArea)
	})

	return func() {
		doInUIThread(func() {
			account.removeCurrentNotificationIf(notification)
		})
	}
}

func (u *gtkUI) notifyTorIsNotRunning(account *account) {
	doInUIThread(func() {
		notification := account.buildTorNotRunningNotification(u)
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func (u *gtkUI) notifyConnectionFailure(account *account, moreInfo func()) {
	doInUIThread(func() {
		notification := account.buildConnectionFailureNotification(u, moreInfo)
		account.setCurrentNotification(notification, u.notificationArea)
	})
}

func buildVerifyIdentityNotification(acc *account, peer, resource string, win gtki.Window) gtki.InfoBar {
	builder := newBuilder("VerifyIdentityNotification")

	obj := builder.getObj("infobar")
	infoBar := obj.(gtki.InfoBar)

	obj = builder.getObj("message")
	message := obj.(gtki.Label)
	message.SetSelectable(true)

	text := fmt.Sprintf(i18n.Local("You have not verified the identity of %s"), peer)
	message.SetText(text)

	obj = builder.getObj("button_verify")
	button := obj.(gtki.Button)
	button.Connect("clicked", func() {
		doInUIThread(func() {
			resp := verifyFingerprintDialog(acc, peer, resource, win)
			if resp == gtki.RESPONSE_YES {
				infoBar.Hide()
				infoBar.Destroy()
			}
		})
	})

	infoBar.ShowAll()

	return infoBar
}

func (u *gtkUI) notify(title, message string) {
	builder := newBuilder("SimpleNotification")
	obj := builder.getObj("dialog")
	dlg := obj.(gtki.MessageDialog)

	dlg.SetProperty("title", title)
	dlg.SetProperty("text", message)
	dlg.SetTransientFor(u.window)

	doInUIThread(func() {
		dlg.Run()
		dlg.Destroy()
	})
}
