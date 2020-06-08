package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type listRoomsView struct {
	accountManager *accountManager
	chatManager    *chatManager

	gtki.Dialog `gtk-widget:"list-chat-rooms"`

	service            gtki.Entry          `gtk-widget:"service"`
	roomsModel         gtki.ListStore      `gtk-widget:"rooms"`
	roomsTreeView      gtki.TreeView       `gtk-widget:"rooms-list-view"`
	roomsTreeContainer gtki.ScrolledWindow `gtk-widget:"room-list-scroll"`
	emptyListLabel     gtki.Label          `gtk-widget:"empty-list-label"`
	account            gtki.ComboBox       `gtk-widget:"accounts"`
	accountsModel      gtki.ListStore      `gtk-widget:"accounts-model"`

	ui *gtkUI
}

func (u *gtkUI) listChatRooms() {
	view := u.newListRoomsView(u.accountManager, u.chatManager)
	view.SetApplication(u.app)
	view.Show()
}

func (u *gtkUI) newListRoomsView(accountManager *accountManager, chatManager *chatManager) gtki.Dialog {
	view := &listRoomsView{
		accountManager: accountManager,
		chatManager:    chatManager,
		ui:             u,
	}

	builder := newBuilder("ListChatRooms")
	err := builder.bindObjects(view)
	if err != nil {
		panic(err)
	}

	builder.ConnectSignals(map[string]interface{}{
		"cancel_handler":             view.Destroy,
		"join_selected_room_handler": view.joinSelectedRoom,
		"fetch_rooms_handler":        view.fetchRoomsFromService,
	})

	doInUIThread(view.populateAccountsModel)

	return view
}

func (v *listRoomsView) fetchRoomsFromService() {
	v.roomsModel.Clear()
	service, _ := v.service.GetText()
	account, found := v.accountManager.getAccountByID(v.getSelectedAccountID())
	if !found {
		return
	}

	conn := account.session.Conn()
	result, _ := conn.GetChatContext().QueryRooms(service)

	doInUIThread(func() {
		if len(result) == 0 {
			v.showLabel(service)
			return
		}

		v.showTreeView()
		for _, room := range result {
			iter := v.roomsModel.Append()
			_ = v.roomsModel.SetValue(iter, 0, room.Name)
			// TODO: parse description?
			_ = v.roomsModel.SetValue(iter, 1, room.Name)
		}
	})
}

func (v *listRoomsView) populateAccountsModel() {
	accs := v.accountManager.getAllConnectedAccounts()
	for _, acc := range accs {
		iter := v.accountsModel.Append()
		_ = v.accountsModel.SetValue(iter, 0, acc.session.GetConfig().Account)
		_ = v.accountsModel.SetValue(iter, 1, acc.session.GetConfig().ID())
	}

	if len(accs) > 0 {
		v.account.SetActive(0)
	}
}

func (v *listRoomsView) showLabel(service string) {
	v.emptyListLabel.SetLabel("No rooms found from service " + service)
	v.emptyListLabel.SetVisible(true)
	v.roomsTreeContainer.SetVisible(false)
}

func (v *listRoomsView) showTreeView() {
	v.emptyListLabel.SetVisible(false)
	v.roomsTreeContainer.SetVisible(true)
}

func (v *listRoomsView) joinSelectedRoom() {
	room := v.getSelectedRoomName()
	service, _ := v.service.GetText()

	addChatView := v.ui.newChatView(v.accountManager, v.chatManager)

	addChatView.service.SetText(service)
	addChatView.room.SetText(room)
	index := v.getSelectedAccountIndex()
	addChatView.setActiveAccount(index)

	v.Destroy()
	addChatView.Show()
}

func (v *listRoomsView) getSelectedRoomName() string {
	ts, _ := v.roomsTreeView.GetSelection()
	_, iter, selected := ts.GetSelected()

	if !selected {
		return ""
	}

	value, _ := v.roomsModel.GetValue(iter, 1)
	roomJid, _ := value.GetString()
	return roomJid
}

func (v *listRoomsView) getSelectedAccountID() string {
	iter, _ := v.account.GetActiveIter()

	val, err := v.accountsModel.GetValue(iter, 1)
	if err != nil {
		return ""
	}

	account, err := val.GetString()
	if err != nil {
		return ""
	}
	return account
}

func (v *listRoomsView) getSelectedAccountIndex() int {
	return v.account.GetActive()
}
