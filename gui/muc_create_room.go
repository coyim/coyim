package gui

import (
	"fmt"

	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	accountManager *accountManager
	errorBox       *errorNotification

	gtki.Dialog `gtk-widget:"create-chat-dialog"`

	notification gtki.Box      `gtk-widget:"notification-area"`
	form         gtki.Grid     `gtk-widget:"form"`
	account      gtki.ComboBox `gtk-widget:"accounts"`
	service      gtki.Entry    `gtk-widget:"service"`
	room         gtki.Entry    `gtk-widget:"room"`

	model       gtki.ListStore `gtk-widget:"accounts-model"`
	accountList []*account

	ui *gtkUI
}

func (u *gtkUI) newMUCRoomView(accountManager *accountManager) *createMUCRoom {
	view := &createMUCRoom{
		accountManager: accountManager,
		ui:             u,
	}

	builder := newBuilder("MUCCreateRoom")
	panicOnDevError(builder.bindObjects(view))
	view.errorBox = newErrorNotification(view.notification)

	view.populateModel(u.getAllConnectedAccounts())

	accountsObserverToken := u.onChangeOfConnectedAccounts(func() {
		doInUIThread(func() {
			view.errorBox.Hide()
			view.populateModel(u.getAllConnectedAccounts())
		})
	})

	builder.ConnectSignals(map[string]interface{}{
		"create_room_handler": view.createRoomHandler,
		"cancel_handler":      view.Destroy,
		"on_close_window_signal": func() {
			u.removeConnectedAccountsObserver(accountsObserverToken)
		},
	})

	return view
}

func (v *createMUCRoom) populateModel(accs []*account) {
	newActiveAccount := 0
	oldActiveAccount := v.account.GetActive()
	if oldActiveAccount >= 0 {
		for i, acc := range accs {
			if acc == v.accountList[oldActiveAccount] {
				newActiveAccount = i
			}
		}
		v.model.Clear()
	}

	for _, acc := range accs {
		iter := v.model.Append()
		_ = v.model.SetValue(iter, 0, acc.session.GetConfig().Account)
		_ = v.model.SetValue(iter, 1, acc.session.GetConfig().ID())
	}

	if len(accs) > 0 {
		v.account.SetActive(newActiveAccount)
	} else {
		v.errorBox.ShowMessage("No accounts connected. Please connect some account from your list of accounts.")
	}
	v.accountList = accs
}

func (v *createMUCRoom) createRoomHandler() {
	idAcc := v.getSelectedAccountID()
	v.errorBox.Hide()

	if idAcc == "" {
		v.errorBox.ShowMessage("No account selected, please select one account from the list or connect some account.")
		return
	}

	account, found := v.accountManager.getAccountByID(idAcc)
	if !found {
		v.errorBox.ShowMessage(fmt.Sprintf("The given account %s is not connected.", idAcc))
		return
	}

	roomName, _ := v.room.GetText()
	service, _ := v.service.GetText()

	if roomName == "" || service == "" {
		v.errorBox.ShowMessage("Please fill the required fields to create the room.")
		return
	}

	complete := make(chan error)
	go func() {
		complete <- account.session.CreateRoom(jid.Parse(fmt.Sprintf("%s@%s", roomName, service)).(jid.Bare))
	}()

	if <-complete != nil {
		v.errorBox.ShowMessage("Could not create the new room")
	} else {
		v.errorBox.ShowMessage("Room created with success")
	}

}

func (v *createMUCRoom) getSelectedAccountID() string {
	iter, _ := v.account.GetActiveIter()

	val, err := v.model.GetValue(iter, 1)
	if err != nil {
		return ""
	}

	account, err := val.GetString()
	if err != nil {
		return ""
	}
	return account
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newMUCRoomView(u.accountManager)
	view.SetApplication(u.app)
	view.Show()
}
