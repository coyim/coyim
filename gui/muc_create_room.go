package gui

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/gotk3adapter/gtki"
)

type createMUCRoom struct {
	accountManager *accountManager
	chatManager    *chatManager
	errorBox       *errorNotification

	gtki.Dialog `gtk-widget:"create-chat-dialog"`

	notification gtki.Box      `gtk-widget:"notification-area"`
	form         gtki.Grid     `gtk-widget:"form"`
	account      gtki.ComboBox `gtk-widget:"accounts"`
	service      gtki.Entry    `gtk-widget:"service"`
	room         gtki.Entry    `gtk-widget:"room"`

	model gtki.ListStore `gtk-widget:"accounts-model"`

	ui *gtkUI
}

func (u *gtkUI) newMUCRoomView(accountManager *accountManager, chatManager *chatManager) *createMUCRoom {
	view := &createMUCRoom{
		accountManager: accountManager,
		chatManager:    chatManager,
		ui:             u,
	}

	builder := newBuilder("MUCCreateRoom")
	panicOnDevError(builder.bindObjects(view))

	builder.ConnectSignals(map[string]interface{}{
		"create_room_handler": view.createRoomHandler,
		"cancel_handler":      view.Destroy,
	})

	view.errorBox = newErrorNotification(view.notification)
	doInUIThread(view.populateModel)

	return view
}

func (v *createMUCRoom) populateModel() {
	accs := v.accountManager.getAllConnectedAccounts()
	for _, acc := range accs {
		iter := v.model.Append()
		_ = v.model.SetValue(iter, 0, acc.session.GetConfig().Account)
		_ = v.model.SetValue(iter, 1, acc.session.GetConfig().ID())
	}

	if len(accs) > 0 {
		v.account.SetActive(0)
	}
}

func (v *createMUCRoom) createRoomHandler() {
	v.errorBox.Hide()
	//TODO: asign creaate room logic
	roomName, _ := v.room.GetText()
	service, _ := v.service.GetText()
	account, found := v.accountManager.getAccountByID(v.getSelectedAccountID())

	if !found {
		return
	}

	conn := account.session.Conn()

	room := &data.LegacyOldDoNotUseRoom{
		ID:      roomName,
		Service: service,
	}
	err := conn.GetChatContext().LegacyOldDoNotUseCreateRoom(room)
	if err != nil {
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

func (u *gtkUI) createChatRoom() {
	view := u.newMUCRoomView(u.accountManager, u.chatManager)
	view.SetApplication(u.app)
	view.Show()
}
