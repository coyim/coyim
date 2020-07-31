package gui

import (
	"fmt"
	"sync"
	"time"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
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
		_ = jrv.modelAccount.SetValue(iter, 0, acc.Account())
		_ = jrv.modelAccount.SetValue(iter, 1, acc.ID())
		jrv.accounts[acc.ID()] = acc
	}

	if len(accounts) > 0 {
		jrv.cmbAccount.SetActive(currentlyActive)
	} else {
		jrv.spinner.Stop()
		jrv.spinner.SetVisible(false)
	}
}

// tryJoinRoom find the room information and make the join to the room
func (u *gtkUI) tryJoinRoom(jrv *mucJoinRoomView, a *account) {
	jrv.updateLock.Lock()

	doInUIThread(jrv.clearErrors)

	roomName, _ := jrv.txtRoomName.GetText()

	doInUIThread(func() {
		jrv.spinner.Start()
		jrv.spinner.SetVisible(true)
	})

	rl, err := a.session.GetRoom(jid.Parse(roomName).(jid.Bare))
	go func() {
		defer func() {
			jrv.updateLock.Unlock()
		}()

		doInUIThread(func() {
			time.Sleep(5000 * time.Millisecond)
			jrv.spinner.Stop()
			jrv.spinner.SetVisible(false)
		})

		if err != nil {
			jrv.notifyOnError(i18n.Local(fmt.Sprintf("The Room \"%s\" doesn't exists", roomName)))
			a.log.Debug(fmt.Sprintf("The Room \"%s\" doesn't exists", roomName))
		} else {
			doInUIThread(func() {
				u.mucShowRoom(a, rl)
				jrv.dialog.Hide()
			})
		}
	}()
}

//
// Custom GTK Events
//

func (jrv *mucJoinRoomView) onShowWindow() {

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
			act := view.cmbAccount.GetActive()
			if act >= 0 && act < len(view.accountsList) && act != view.currentlyActive {
				view.currentlyActive = act
			}
		},
		"on_btn_cancel_clicked_signal": view.dialog.Destroy,
		"on_btn_join_clicked_signal": func() {
			idx := view.cmbAccount.GetActive()
			act := view.accountsList[idx]
			u.tryJoinRoom(view, act)
		},
	})

	u.connectShortcutsChildWindow(view.dialog)

	view.dialog.SetTransientFor(u.window)
	view.dialog.Show()
}
