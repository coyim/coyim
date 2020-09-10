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

	identity jid.Bare
	occupant string
	info     *muc.RoomListing

	log      coylog.Logger
	joined   bool
	returnTo func()

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

func (u *gtkUI) newRoomView(a *account, ident jid.Bare, roomInfo *muc.RoomListing) *roomView {
	view := &roomView{
		u:        u,
		account:  a,
		identity: ident,
		info:     roomInfo,
	}

	view.log = a.log.WithField("room", ident)

	view.initUIBuilder()
	view.initDefaults()

	toolbar := newRoomViewToolbar()
	view.toolbar = toolbar

	roster := newRoomViewRoster()
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
	room.Opaque = u.newRoomView(a, ident, roomInfo)
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

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare, roomInfo *muc.RoomListing, returnTo func()) {
	room, ok := u.getRoomOrCreateItIfNoExists(a, ident, roomInfo)

	view := getViewFromRoom(room)
	if !view.joined {
		view.returnTo = returnTo

		view.switchToLobbyView(view.info)
	} else {
		// In the main view of the room, we don't have the "cancel"
		// functionality that it's useful only in the lobby view of the room.
		// For that reason is why we ignore the "returnTo" value.
		view.returnTo = nil

		view.switchToMainView()
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
	go v.forceLeaveRoom()
}

func (v *roomView) forceLeaveRoom() {
	if !v.joined {
		v.log.Info("The occupant left the lobby")
		v.leaveRoomMananger()
		return
	}
	// TODO: Should we implement a timeout here?
	for {
		sc := make(chan bool)
		go func() {
			v.tryLeaveRoom(func() {
				sc <- true
			}, func() {
				sc <- false
			})
		}()

		if <-sc {
			return
		}
	}
}

func (v *roomView) tryLeaveRoom(s, f func()) {
	go func() {
		resultCh, errCh := v.account.session.LeaveRoom(v.identity, v.occupant)
		select {
		case <-resultCh:
			v.onLeaveRoomSuccess(s)
		case err := <-errCh:
			v.log.WithError(err).Error("An error occurred when trying to leave the room")
			v.onLeaveRoomFailure(f)
		}
	}()
}

func (v *roomView) onLeaveRoomSuccess(f func()) {
	v.log.Info("The occupant left the room")
	v.leaveRoomMananger()
	v.joined = false
	if f != nil {
		f()
	}
}

func (v *roomView) leaveRoomMananger() {
	v.account.roomManager.LeaveRoom(v.identity)
}

func (v *roomView) onLeaveRoomFailure(f func()) {
	if f != nil {
		f()
	}
}

func (v *roomView) switchToLobbyView(roomInfo *muc.RoomListing) {
	if v.lobby == nil {
		v.lobby = newRoomViewLobby(v.account, v.identity, v.content, v.onEntered, v.onCancel, roomInfo)
	} else {
		v.lobby.setRoomInfo(roomInfo)
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

func (v *roomView) onEntered() {
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

// onRoomOccupantJoinedReceived SHOULD be called from the UI thread
func (v *roomView) onRoomOccupantJoinedReceived(occupant string, occupants []*muc.Occupant) {
	if v.joined {
		v.log.WithField("occupant", occupant).Error("A joined event was received but the user is already in the room")
		return
	}

	v.occupant = occupant
	v.lobby.onRoomOccupantJoinedReceived()
	v.roster.updateRoomRoster(occupants)
}

// onRoomOccupantUpdateReceived SHOULD be called from the UI thread
func (v *roomView) onRoomOccupantUpdateReceived(occupants []*muc.Occupant) {
	v.roster.updateRoomRoster(occupants)
}

// onRoomOccupantLeftTheRoomReceived SHOULD be called from the UI thread
func (v *roomView) onRoomOccupantLeftTheRoomReceived(occupant jid.Resource, occupants []*muc.Occupant) {
	v.conv.showOccupantLeftRoom(occupant)
	v.roster.updateRoomRoster(occupants)
}

// onRoomMessageToTheRoomReceived should be called from the UI thread
func (v *roomView) onRoomMessageToTheRoomReceived(occupant jid.Resource, message string) {
	v.conv.showMessageInChatRoom(occupant, message, mtLiveMessage)
}

// loggingIsEnabled MUST not be called from the UI thread
func (v *roomView) loggingIsEnabled() {
	if v.conv != nil {
		msg := i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website.")
		doInUIThread(func() {
			v.conv.showMessageInChatRoom(jid.Resource{}, msg, mtWarning)
		})
	}
}

// loggingIsDisabled MUST not be called from the UI thread
func (v *roomView) loggingIsDisabled() {
	if v.conv != nil {
		msg := i18n.Local("This room is no longer publicly logged.")
		doInUIThread(func() {
			v.conv.showMessageInChatRoom(jid.Resource{}, msg, mtWarning)
		})
	}
}
