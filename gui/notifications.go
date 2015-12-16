package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"../i18n"
)

func (u *gtkUI) showConnectAccountNotification(account *account) {
	account.buildConnectionNotification()

	glib.IdleAdd(func() {
		infoBar := account.connectionNotification
		u.notificationArea.Add(infoBar)
		infoBar.ShowAll()
	})
}

func (u *gtkUI) removeConnectAccountNotification(account *account) {
	glib.IdleAdd(func() {
		account.removeConnectionNotification()
	})
}

func (u *gtkUI) notifyConnectionFailure(account *account) {
	builder := builderForDefinition("ConnectionFailureNotification")

	builder.ConnectSignals(map[string]interface{}{
		"handleResponse": func(info *gtk.InfoBar, response gtk.ResponseType) {
			if response != gtk.RESPONSE_CLOSE {
				return
			}

			info.Hide()
			info.Destroy()
		},
	})

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	message := obj.(*gtk.Label)

	text := fmt.Sprintf(i18n.Local("Connection failure\n%s"),
		account.session.CurrentAccount.Account)
	message.SetText(text)

	glib.IdleAdd(func() {
		u.notificationArea.Add(infoBar)
		infoBar.ShowAll()
	})
}

func buildVerifyIdentityNotification(peer string) *gtk.InfoBar {
	builder := builderForDefinition("VerifyIdentityNotification")

	obj, _ := builder.GetObject("infobar")
	infoBar := obj.(*gtk.InfoBar)

	obj, _ = builder.GetObject("message")
	message := obj.(*gtk.Label)

	text := fmt.Sprintf(i18n.Local("You have not verified the identity of %s"), peer)
	message.SetText(text)

	return infoBar
}
