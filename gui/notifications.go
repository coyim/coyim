package gui

import "github.com/gotk3/gotk3/glib"

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
