package gui

import (
	"sync"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucJoinRoomView struct {
	builder *builder

	generation int
	updateLock sync.RWMutex

	dialog           gtki.Dialog    `gtk-widget:"MUCJoinRoom"`
	cmbAccount       gtki.ComboBox  `gtk-widget:"cmbAccount"`
	txtRoomName      gtki.Entry     `gtk-widget:"txtRoomName"`
	spinner          gtki.Spinner   `gtk-widget:"spinner"`
	notificationArea gtki.Box       `gtk-widget:"boxNotificationArea"`
	modelAccount     gtki.ListStore `gtk-widget:"modelAccount"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification

	accountsList    []*account
	accounts        map[string]*account
	currentlyActive int
}

func (jrv *mucJoinRoomView) clearErrors() {
	jrv.errorNotif.Hide()
}

func (jrv *mucJoinRoomView) notifyOnError(errMessage string) {
	doInUIThread(func() {
		if jrv.notification != nil {
			jrv.notificationArea.Remove(jrv.notification)
		}

		jrv.errorNotif.ShowMessage(errMessage)
	})
}

func (jrv *mucJoinRoomView) init() {
	jrv.builder = newBuilder("MUCJoinRoomDialog")
	panicOnDevError(jrv.builder.bindObjects(jrv))
	jrv.errorNotif = newErrorNotification(jrv.notificationArea)
}

// initOrReplaceAccounts should be called from the UI thread
func (jrv *mucJoinRoomView) initOrReplaceAccounts(accounts []*account) {
	if len(accounts) == 0 {
		jrv.notifyOnError(i18n.Local("There are no connected accounts"))
	}

	currentlyActive := 0
	oldActive := jrv.cmbAccount.GetActive()
	if jrv.accounts != nil && oldActive >= 0 {
		currentlyActiveAccount := jrv.accountsList[oldActive]
		for ix, a := range accounts {
			if currentlyActiveAccount == a {
				currentlyActive = ix
				jrv.currentlyActive = currentlyActive
			}
		}
		jrv.modelAccount.Clear()
	}

	jrv.accountsList = accounts
	jrv.accounts = make(map[string]*account)
	for _, acc := range accounts {
		iter := jrv.modelAccount.Append()
		_ = jrv.modelAccount.SetValue(iter, 0, acc.session.GetConfig().Account)
		_ = jrv.modelAccount.SetValue(iter, 1, acc.session.GetConfig().ID())
		jrv.accounts[acc.session.GetConfig().ID()] = acc
	}

	if len(accounts) > 0 {
		jrv.cmbAccount.SetActive(currentlyActive)
	} else {
		jrv.spinner.Stop()
		jrv.spinner.SetVisible(false)
	}
}

func (jrv *mucJoinRoomView) onShowWindow() {

}

func (jrv *mucJoinRoomView) onCmbAccountChanged() {
	act := jrv.cmbAccount.GetActive()
	if act >= 0 && act < len(jrv.accountsList) && act != jrv.currentlyActive {
		jrv.currentlyActive = act
	}
}

func (jrv *mucJoinRoomView) onBtnJoinClicked() {

}

// mucJoinRoom should be called from the UI thread
func (u *gtkUI) mucShowJoinRoom() {
	view := &mucJoinRoomView{}

	view.init()

	view.initOrReplaceAccounts(u.getAllConnectedAccounts())

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal": func() {},
		"on_show_window_signal": func() {
			view.onShowWindow()
		},
		"on_cmb_account_changed": func() {
			view.onCmbAccountChanged()
		},
		"on_btn_cancel_clicked_signal": view.dialog.Destroy,
		"on_btn_join_clicked_signal": func() {
			view.onBtnJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
