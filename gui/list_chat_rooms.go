package gui

import (
	"github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/interfaces"
	"github.com/coyim/gotk3adapter/gtki"
)

type listRoomsView struct {
	accountManager *accountManager
	chatManager    *chatManager
	errorBox       *errorNotification

	gtki.Dialog `gtk-widget:"list-chat-rooms"`

	service            gtki.Entry          `gtk-widget:"service"`
	roomsModel         gtki.ListStore      `gtk-widget:"rooms"`
	roomsTreeView      gtki.TreeView       `gtk-widget:"rooms-list-view"`
	roomsTreeContainer gtki.ScrolledWindow `gtk-widget:"room-list-scroll"`
	emptyListLabel     gtki.Label          `gtk-widget:"empty-list-label"`
}

func (u *gtkUI) listChatRooms() {
	view := newListRoomsView(u.accountManager, u.chatManager)
	view.SetTransientFor(u.window)
	view.Show()
}

func newListRoomsView(accountManager *accountManager, chatManager *chatManager) gtki.Dialog {
	view := &listRoomsView{
		accountManager: accountManager,
		chatManager:    chatManager,
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

	return view
}

func (v *listRoomsView) fetchRoomsFromService() {
	v.roomsModel.Clear()
	service, _ := v.service.GetText()

	//TODO: Be able to select account
	account := v.accountManager.getAllAccounts()[0]

	conn := account.session.Conn()
	result, _ := conn.GetChatContext().QueryRooms(service)

	doInUIThread(func() {
		if len(result) == 0 {
			v.emptyListLabel.SetLabel("No rooms found from service " + service)
			v.showLabel()
			return
		}

		v.showTreeView()
		for _, room := range result {
			iter := v.roomsModel.Append()
			v.roomsModel.SetValue(iter, 0, room.Jid)
			v.roomsModel.SetValue(iter, 1, room.Name)
			v.roomsModel.SetValue(iter, 2, room.Name)
		}
	})
}

func (v *listRoomsView) showLabel() {
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

	addChatView := newChatView(v.accountManager, v.chatManager)
	if parent, err := v.GetTransientFor(); err == nil {
		addChatView.SetTransientFor(parent)
	}

	addChatView.service.SetText(service)
	addChatView.room.SetText(room)

	v.Destroy()
	addChatView.Show()
}

func (v *listRoomsView) getChatContext(eventsChan chan interface{}) interfaces.Chat {
	chat, err := v.chatManager.getChatContextForAccount(v.getHandle(), eventsChan)
	if err != nil {
		v.errorBox.ShowMessage(err.Error())
		return nil
	}
	return chat
}

func (v *listRoomsView) getSelectedRoomName() string {
	ts, _ := v.roomsTreeView.GetSelection()
	_, iter, selected := ts.GetSelected()

	if !selected {
		//TODO: Error handling
		return ""
	}

	value, _ := v.roomsModel.GetValue(iter, 1)
	roomJid, _ := value.GetString()
	return roomJid
}

func (v *listRoomsView) getHandle() string {
	return v.accountManager.getAllAccounts()[0].ID()
}

func (v *listRoomsView) buildOccupant(room, service, handle string) *data.Occupant {

	return &data.Occupant{
		Room: data.Room{
			ID:      room,
			Service: service,
		},
		Handle: handle,
	}
}
