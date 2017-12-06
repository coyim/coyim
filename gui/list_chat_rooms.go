package gui

import (
	"log"

	"github.com/coyim/gotk3adapter/gtki"
)

type listRoomsView struct {
	accountManager *accountManager

	gtki.Dialog `gtk-widget:"list-chat-rooms"`

	service       gtki.Entry     `gtk-widget:"service"`
	roomsModel    gtki.ListStore `gtk-widget:"rooms"`
	roomsTreeView gtki.TreeView  `gtk-widget:"rooms-list-view"`
}

func (u *gtkUI) listChatRooms() {
	view := newListRoomsView(u.accountManager)
	view.SetTransientFor(u.window)
	view.Show()
}

func newListRoomsView(accountManager *accountManager) gtki.Dialog {
	view := &listRoomsView{
		accountManager: accountManager,
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
	//TODO: deal with empty results

	doInUIThread(func() {
		for _, room := range result {
			iter := v.roomsModel.Append()
			v.roomsModel.SetValue(iter, 0, room.Jid)
			v.roomsModel.SetValue(iter, 1, room.Name)
			v.roomsModel.SetValue(iter, 2, room.Name)
		}
	})
}

func (v *listRoomsView) joinSelectedRoom() {
	ts, _ := v.roomsTreeView.GetSelection()
	if _, iter, ok := ts.GetSelected(); ok {
		value, _ := v.roomsModel.GetValue(iter, 0)
		roomJid, _ := value.GetString()
		log.Print("ROOM: ")
		log.Print(roomJid)
	}
}
