package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/i18n"
)

func (u *gtkUI) showConnectAccountNotification(account *account) {
	err := account.buildConnectionNotification()
	if err != nil {
		panic(err)
	}

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
	builder, err := loadBuilderWith("ConnectionFailureNotification")
	if err != nil {
		panic(err)
	}

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

	text := fmt.Sprintf(i18n.Local("Connection lost\n%s"),
		account.session.CurrentAccount.Account)
	message.SetText(text)

	glib.IdleAdd(func() {
		u.notificationArea.Add(infoBar)
		infoBar.ShowAll()
	})
}
