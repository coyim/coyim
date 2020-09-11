package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	log "github.com/sirupsen/logrus"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomView struct {
	u       *gtkUI
	account *account
	builder *builder

	identity   jid.Bare
	occupant   *muc.Occupant
	roomRoster *muc.RoomRoster
	info       *muc.RoomListing

	log      coylog.Logger
	joined   bool
	returnTo func()

	window  gtki.Window `gtk-widget:"roomWindow"`
	content gtki.Box    `gtk-widget:"boxMainView"`

	onSelfJoinedReceived     []func()
	onOccupantJoinedReceived []func()

	main    *roomViewMain
	toolbar *roomViewToolbar
	roster  *roomViewRoster
	conv    *roomViewConversation
	lobby   *roomViewLobby
}

func getViewFromRoom(r *muc.Room) *roomView {
	return r.Opaque.(*roomView)
}

func (v *roomView) onSelfJoinReceived(f func()) {
	v.onSelfJoinedReceived = append(v.onSelfJoinedReceived, f)
}

func (v *roomView) onOccupantReceived(f func()) {
	v.onOccupantJoinedReceived = append(v.onOccupantJoinedReceived, f)
}

func (u *gtkUI) newRoomView(a *account, ident jid.Bare, roomInfo *muc.RoomListing, roomRoster *muc.RoomRoster) *roomView {
	view := &roomView{
		u:          u,
		account:    a,
		identity:   ident,
		info:       roomInfo,
		roomRoster: roomRoster,
	}

	view.log = a.log.WithField("room", ident)

	view.initUIBuilder()
	view.initDefaults()

	toolbar := view.newRoomViewToolbar()
	view.toolbar = toolbar

	roster := view.newRoomViewRoster()
	view.roster = roster

	conversation := u.newRoomViewConversation()
	view.conv = conversation

	return view
}

func (v *roomView) setTitle(r string) {
	v.window.SetTitle(r)
}

func (u *gtkUI) newRoom(a *account, ident jid.Bare, roomInfo *muc.RoomListing) *muc.Room {
	room := muc.NewRoom(ident)
	room.Opaque = u.newRoomView(a, ident, roomInfo, room.Roster())
	return room
}

func (u *gtkUI) getRoomOrCreateItIfNoExists(a *account, ident jid.Bare, roomInfo *muc.RoomListing) (*muc.Room, bool) {
	room, ok := a.roomManager.GetRoom(ident)
	if !ok {
		room = u.newRoom(a, ident, roomInfo)
		a.roomManager.AddRoom(room)
	}
	return room, ok
}

// mucShowRoom MUST be called always from the UI thread
//
// Also, when we want to show a chat room, having a "return to" function that
// will be called from the lobby only when the user wants to "cancel" or "return"
// might be useful in some scenarios like "returning to previous step".
//
// Please note that "returnTo" will be called from the UI thread too
func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare, roomInfo *muc.RoomListing, returnTo func()) {
	room, ok := u.getRoomOrCreateItIfNoExists(a, ident, roomInfo)

	view := getViewFromRoom(room)
	if view.joined {
		// In the main view of the room, we don't have the "cancel"
		// functionality that it's useful only in the lobby view of the room.
		// For that reason is why we ignore the "returnTo" value.
		view.returnTo = nil

		view.switchToMainView()
	} else {
		view.returnTo = returnTo
		view.switchToLobbyView(view.info)
	}

	if !ok {
		view.window.Show()
	} else {
		view.window.Present()
	}
}

func (v *roomView) initUIBuilder() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_destroy_window": v.onDestroyWindow,
	})
}

func (v *roomView) initDefaults() {
	v.setTitle(v.identity.String())
}

func (v *roomView) onDestroyWindow() {
	v.leaveRoomMananger()
}

func (v *roomView) LeaveRoom() {
	go func() {
		resultCh, errCh := v.account.session.LeaveRoom(v.identity, v.occupant.Nick)
		select {
		case <-resultCh:
			v.leaveRoomMananger()
			doInUIThread(v.window.Destroy)
		case err := <-errCh:
			v.log.WithError(err).Error("An error occurred when trying to leave the room")
			//TODO: Show an appropiate way to present the error message to the user.
		}
	}()
}

func (v *roomView) leaveRoomMananger() {
	v.account.roomManager.LeaveRoom(v.identity)
}

func (v *roomView) switchToLobbyView(roomInfo *muc.RoomListing) {
	if v.lobby == nil {
		v.lobby = v.newRoomViewLobby(v.account, v.identity, v.content, v.onJoined, v.onJoinCancel, roomInfo)
	} else {
		// If we got new room information, we should show
		// any warnings based on that info
		v.lobby.showRoomWarnings(roomInfo)
	}

	if v.shouldReturnOnCancel() {
		v.lobby.swtichToReturnOnCancel()
	} else {
		v.lobby.swtichToCancel()
	}

	v.lobby.show()
}

func (v *roomView) shouldReturnOnCancel() bool {
	return v.returnTo != nil
}

func (v *roomView) switchToMainView() {
	if v.main == nil {
		v.main = newRoomMainView(v.account, v.identity, v.conv.view, v.roster.view, v.toolbar.view, v.content)
	}
	v.main.show()
}

func (v *roomView) onJoined() {
	v.joined = true

	doInUIThread(func() {
		v.lobby.hide()
		v.switchToMainView()
	})
}

// TODO: if we have an active connection or request, we should
// stop/close it here before destroying the window
func (v *roomView) onJoinCancel() {
	v.window.Destroy()

	if v.shouldReturnOnCancel() {
		v.returnTo()
	}
}

func (v *roomView) onNicknameConflictReceived(room jid.Bare, nickname string) {
	if v.joined {
		v.log.WithFields(log.Fields{
			"room":     room,
			"nickname": nickname,
		}).Error("A nickname conflict event was received but the user is already in the room")
		return
	}

	v.lobby.onNicknameConflictReceived(room, nickname)
}

func (v *roomView) onRegistrationRequiredReceived(room jid.Bare, nickname string) {
	if v.joined {
		v.log.WithFields(log.Fields{
			"room":     room,
			"nickname": nickname,
		}).Error("A registration required event was received but the user is already in the room")
		return
	}

	v.lobby.onRegistrationRequiredReceived(room, nickname)
}

func (v *roomView) onRoomOccupantErrorReceived(room jid.Bare, nickname string) {
	if v.joined {
		v.log.WithFields(log.Fields{
			"room":     room,
			"nickname": nickname,
		}).Error("A joined event error was received but the user is already in the room")
		return
	}

	v.lobby.onJoinErrorRecevied(room, nickname)
}

// onRoomOccupantJoinedReceived MUST be called from the UI thread
func (v *roomView) onRoomOccupantJoinedReceived(occupant *muc.Occupant) {
	if v.joined {
		v.log.WithField("occupant", occupant).Error("A joined event was received but the user is already in the room")
		return
	}

	v.occupant = occupant
	for _, f := range v.onSelfJoinedReceived {
		f()
	}
}

// onRoomOccupantUpdateReceived MUST be called from the UI thread
func (v *roomView) onRoomOccupantUpdateReceived() {
	// TODO[OB] - This is incorrect
	for _, f := range v.onOccupantJoinedReceived {
		f()
	}
}

// onRoomOccupantLeftTheRoomReceived MUST be called from the UI thread
func (v *roomView) onRoomOccupantLeftTheRoomReceived(occupant jid.Resource) {
	if v.conv != nil {
		v.conv.showOccupantLeftRoom(occupant.String())
	}
	v.roster.updateRosterModel()
}

// someoneJoinedTheRoom MUST be called from the UI thread
func (v *roomView) someoneJoinedTheRoom(nick string) {
	if v.conv != nil {
		v.conv.displayNotificationWhenOccupantJoinedRoom(nick)
	}
	v.roster.updateRosterModel()
}

// onRoomMessageToTheRoomReceived MUST be called from the UI thread
func (v *roomView) onRoomMessageToTheRoomReceived(nickname, subject, message string) {
	if v.conv != nil {
		v.conv.showLiveMessageInTheRoom(nickname, subject, message)
	}
}

// loggingIsEnabled MUST not be called from the UI thread
func (v *roomView) loggingIsEnabled() {
	if v.conv != nil {
		msg := i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website.")
		doInUIThread(func() {
			v.conv.addLineToChatTextUsingTagID(msg, "warning")
		})
	}
}

// loggingIsDisabled MUST not be called from the UI thread
func (v *roomView) loggingIsDisabled() {
	if v.conv != nil {
		msg := i18n.Local("This room is no longer publicly logged.")
		doInUIThread(func() {
			v.conv.addLineToChatTextUsingTagID(msg, "warning")
		})
	}
}
