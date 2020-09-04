package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomView struct {
	u       *gtkUI
	account *account
	builder *builder

	log      coylog.Logger
	identity jid.Bare
	joined   bool

	window  gtki.Window `gtk-widget:"roomWindow"`
	content gtki.Box    `gtk-widget:"boxMainView"`

	main    *roomViewMain
	toolbar *roomViewToolbar
	roster  *roomViewRoster
	conv    *roomViewConversation
	lobby   *roomViewLobby
}

func getViewFromRoom(r *muc.Room) *roomView {
	return r.Opaque.(*roomView)
}

func (u *gtkUI) newRoomView(a *account, ident jid.Bare) *roomView {
	view := &roomView{
		u:        u,
		account:  a,
		identity: ident,
	}

	view.initUIBuilder()
	view.initDefaults()

	toolbar := newRoomViewToolbar()
	view.toolbar = toolbar

	roster := newRoomViewRoster()
	view.roster = roster

	conversation := newRoomViewConversation()
	view.conv = conversation

	return view
}

func (v *roomView) setTitle(r string) {
	v.window.SetTitle(r)
}

func (u *gtkUI) newRoom(a *account, ident jid.Bare) *muc.Room {
	room := muc.NewRoom(ident)
	room.Opaque = u.newRoomView(a, ident)
	return room
}

func (u *gtkUI) getRoomOrCreateItIfNoExists(a *account, ident jid.Bare) (*muc.Room, bool) {
	room, ok := a.roomManager.GetRoom(ident)
	if !ok {
		room = u.newRoom(a, ident)
		a.roomManager.AddRoom(room)
	}
	return room, ok
}

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare) {
	room, wasCreated := u.getRoomOrCreateItIfNoExists(a, ident)
	view := getViewFromRoom(room)

	if !wasCreated {
		view.switchToLobbyView()
		view.window.Show()
		return
	}

	view.switchToMainView()
}

func (v *roomView) initUIBuilder() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_show_window":  v.onShowWindow,
		"on_close_window": v.onCloseWindow,
	})
}

func (v *roomView) initDefaults() {
	v.setTitle(v.identity.String())
}

func (v *roomView) onShowWindow() {}

func (v *roomView) onCloseWindow() {
	exists := v.account.roomManager.LeaveRoom(v.identity)
	if !exists {
		v.log.Error("Trying to leave a room that doesn't exists.")
		return
	}
	v.joined = false
}

func (v *roomView) switchToLobbyView() {
	if v.joined {
		panic("developer error: the user is already in this room")
	}

	v.lobby = newRoomViewLobby(v.account, v.identity, v.content, v.onEnter, v.onCancel)
	v.lobby.show()
}

func (v *roomView) switchToMainView() {
	if !v.joined {
		panic("developer error: the user is not in this room")
	}

	v.main = newRoomMainView(v.account, v.identity, v.conv.view, v.roster.view, v.toolbar.view, v.content)
	v.main.show()
}

func (v *roomView) onEnter() {
	v.joined = true
	doInUIThread(func() {
		v.lobby.hide()
		v.switchToMainView()
	})
}

// TODO: if we have an active connection or request, we should
// stop/close it here before destroying the window
func (v *roomView) onCancel() {
	doInUIThread(v.window.Destroy)
}

func (v *roomView) onNicknameConflictReceived(from jid.Full) {
	if v.joined {
		v.log.WithField("from", from).Error("A nickname conflict event was received but the user is already in the room")
		return
	}
	v.lobby.onNicknameConflictReceived(from)
}

func (v *roomView) onRegistrationRequiredReceived(from jid.Full) {
	if v.joined {
		v.log.WithField("from", from).Error("A registration required event was received but the user is already in the room")
		return
	}
	v.lobby.onRegistrationRequiredReceived(from)
}

func (v *roomView) onRoomOccupantErrorReceived(from jid.Full) {
	if v.joined {
		v.log.WithField("from", from).Error("A joined event error was received but the user is already in the room")
		return
	}
	v.lobby.onJoinErrorRecevied(from)
}

// onRoomOccupantJoinedReceived should be called from the UI thread
func (v *roomView) onRoomOccupantJoinedReceived(occupant jid.Resource, occupants []*muc.Occupant) {
	if v.joined {
		v.log.WithField("occupant", occupant).Error("A joined event was received but the user is already in the room")
		return
	}
	v.lobby.onRoomOccupantJoinedReceived()
	v.roster.updateRoomRoster(occupants)
}

// onRoomOccupantUpdateReceived should be called from the UI thread
func (v *roomView) onRoomOccupantUpdateReceived(occupants []*muc.Occupant) {
	v.roster.updateRoomRoster(occupants)
}

// onRoomOccupantLeftTheRoomReceived should be called from the UI thread
func (v *roomView) onRoomOccupantLeftTheRoomReceived(occupant jid.Resource, occupants []*muc.Occupant) {
	v.conv.showOccupantLeftRoom(occupant)
	v.roster.updateRoomRoster(occupants)
}
