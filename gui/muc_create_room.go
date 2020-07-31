package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	accountManager *accountManager
	errorBox       *errorNotification

	gtki.Dialog `gtk-widget:"create-chat-dialog"`

	notification     gtki.Box          `gtk-widget:"notification-area"`
	form             gtki.Grid         `gtk-widget:"form"`
	account          gtki.ComboBox     `gtk-widget:"accounts"`
	chatServices     gtki.ComboBoxText `gtk-widget:"chatServices"`
	chatServiceEntry gtki.Entry        `gtk-widget:"chatServiceEntry"`
	room             gtki.Entry        `gtk-widget:"room"`
	createButton     gtki.Button       `gtk-widget:"button-ok"`
	cancelButton     gtki.Button       `gtk-widget:"button-cancel"`

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

	accountsObserverToken := u.onChangeOfConnectedAccounts(func() {
		doInUIThread(func() {
			view.errorBox.Hide()
			view.populateModel(u.getAllConnectedAccounts())
			view.checkIfFieldsAreEmpty()
		})
	})

	builder.ConnectSignals(map[string]interface{}{
		"create_room_handler": view.createRoomHandler,
		"cancel_handler":      view.Destroy,
		"on_close_window_signal": func() {
			u.removeConnectedAccountsObserver(accountsObserverToken)
		},
		"changed_value_listener":      view.updateChatServices,
		"on_room_changed":             view.checkIfFieldsAreEmpty,
		"on_chatServiceEntry_changed": view.checkIfFieldsAreEmpty,
	})

	view.populateModel(u.getAllConnectedAccounts())

	return view
}

func (v *createMUCRoom) updateChatServices() {
	enteredService, _ := v.chatServiceEntry.GetText()
	v.clearCurrentChatServices()

	acc := v.getCurrentConnectedAcount()
	if acc == nil {
		return
	}

	items, err := acc.session.GetChatServices(jid.Parse(acc.Account()).Host())
	if err != nil {
		return
	}

	for _, i := range items {
		v.chatServices.AppendText(i.Jid)
	}

	if enteredService != "" {
		v.chatServiceEntry.SetText(enteredService)
	} else {
		v.chatServices.SetActive(0)
	}
}

func (v *createMUCRoom) clearCurrentChatServices() {
	v.chatServices.RemoveAll()
	v.chatServiceEntry.SetText("")
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
	} else {
		v.errorBox.ShowMessage(i18n.Local("No accounts connected. Please connect some account from your list of accounts."))
	}
	v.accountList = accs
}

func (v *createMUCRoom) updateFields(f bool) {
	v.cancelButton.SetSensitive(f)
	v.createButton.SetSensitive(f)
	v.account.SetSensitive(f)
	v.room.SetSensitive(f)
	v.chatServices.SetSensitive(f)
}

func (v *createMUCRoom) createRoomHandler() {
	account := v.getCurrentConnectedAcount()
	if account == nil {
		return
	}

	roomName, _ := v.room.GetText()
	service := v.chatServices.GetActiveText()

	if roomName == "" || service == "" {
		v.errorBox.ShowMessage(i18n.Local("Please fill the required fields to create the room."))
		return
	}

	v.updateFields(false)
	originalLabel, _ := v.createButton.GetProperty("label")
	v.createButton.SetProperty("label", i18n.Local("Creating room..."))

	ec := make(chan error)

	go func() {
		ec <- account.session.CreateRoom(jid.Parse(fmt.Sprintf("%s@%s", roomName, service)).(jid.Bare))
	}()

	go func() {
		err, ok := <-ec
		if !ok || err != nil {
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

func (v *createMUCRoom) getCurrentConnectedAcount() *account {
	v.errorBox.Hide()
	idAcc := v.getSelectedAccountID()
	if idAcc == "" {
		v.errorBox.ShowMessage(i18n.Local("No account selected, please select one account from the list or connect to one."))
		return nil
	}

	account, found := v.accountManager.getAccountByID(idAcc)
	if !found {
		v.errorBox.ShowMessage(i18n.Localf("The given account %s is not connected.", idAcc))
		return nil
	}

	return account
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

func (v *createMUCRoom) checkIfFieldsAreEmpty() {
	accountVal := v.getSelectedAccountID()
	serviceVal := v.chatServices.GetActiveText()
	roomVal, _ := v.room.GetText()

	if accountVal == "" || serviceVal == "" || roomVal == "" {
		v.createButton.SetSensitive(false)
	} else {
		v.createButton.SetSensitive(true)
	}
}

func (u *gtkUI) mucCreateChatRoom() {
	view := u.newMUCRoomView(u.accountManager)
	view.SetTransientFor(u.window)
	doInUIThread(view.Show)
}
