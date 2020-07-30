package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

var createMUCRoomIsOpen bool

type createMUCRoom struct {
	accountManager *accountManager
	errorBox       *errorNotification

	gtki.Dialog `gtk-widget:"create-chat-dialog"`

	notification gtki.Box      `gtk-widget:"notification-area"`
	form         gtki.Grid     `gtk-widget:"form"`
	account      gtki.ComboBox `gtk-widget:"accounts"`
	service      gtki.Entry    `gtk-widget:"service"`
	room         gtki.Entry    `gtk-widget:"room"`
	createButton gtki.Button   `gtk-widget:"button-ok"`
	cancelButton gtki.Button   `gtk-widget:"button-cancel"`

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
			createMUCRoomIsOpen = false
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
		_ = v.model.SetValue(iter, 0, acc.Account())
		_ = v.model.SetValue(iter, 1, acc.ID())
	}

	if len(accs) > 0 {
		v.account.SetActive(newActiveAccount)
		v.createButton.SetSensitive(true)
	} else {
		v.errorBox.ShowMessage(i18n.Local("No accounts connected. Please connect some account from your list of accounts."))
		v.createButton.SetSensitive(false)
	}
	v.accountList = accs
}

func (v *createMUCRoom) updateFields(f bool) {
	v.cancelButton.SetSensitive(f)
	v.createButton.SetSensitive(f)
	v.account.SetSensitive(f)
	v.room.SetSensitive(f)
	v.service.SetSensitive(f)
}

func (v *createMUCRoom) createRoomHandler() {
	idAcc := v.getSelectedAccountID()

	v.errorBox.Hide()

	if idAcc == "" {
		v.errorBox.ShowMessage(i18n.Local("No account selected, please select one account from the list or connect some account."))
		return
	}

	account, found := v.accountManager.getAccountByID(idAcc)
	if !found {
		v.errorBox.ShowMessage(i18n.Localf("The given account %s is not connected.", idAcc))
		return
	}

	roomName, _ := v.room.GetText()
	service, _ := v.service.GetText()

	if roomName == "" || service == "" {
		v.errorBox.ShowMessage(i18n.Local("Please fill the required fields to create the room."))
		return
	}

	v.updateFields(false)
	originalLabel, _ := v.createButton.GetProperty("label")
	v.createButton.SetProperty("label", i18n.Local("Creating room..."))

	complete := make(chan error)

	go func() {
		complete <- account.session.CreateRoom(jid.Parse(fmt.Sprintf("%s@%s", roomName, service)).(jid.Bare))
	}()

	go func() {
		if <-complete != nil {
			v.errorBox.ShowMessage(i18n.Local("Could not create the new room"))
		} else {
			v.errorBox.ShowMessage(i18n.Local("Room created with success"))
		}

		doInUIThread(func() {
			v.updateFields(true)
			v.createButton.SetProperty("label", originalLabel)
		})
	}()
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
	if createMUCRoomIsOpen {
		return
	}
	view := u.newMUCRoomView(u.accountManager)
	doInUIThread(view.Show)
	createMUCRoomIsOpen = true
}
