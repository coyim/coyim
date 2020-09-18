package gui

import (
	"sync"

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
	room     *muc.Room
	occupant *muc.Occupant
	info     *muc.RoomListing

	log      coylog.Logger
	joined   bool
	opened   bool
	returnTo func()

	window           gtki.Window  `gtk-widget:"roomWindow"`
	content          gtki.Box     `gtk-widget:"boxMainView"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"roomNotificationArea"`

	notification gtki.InfoBar
	errorNotif   *errorNotification

	// TODO: It's neccessary handle events signals in a better way.
	// Maybe we can take inspiration from the session event management.
	selfJoinedReceived     []func()
	occupantReceived       []func()
	roomInfoReceived       []func(*muc.RoomListing)
	selfJoinedReceivedLock sync.RWMutex
	occupantReceivedLock   sync.RWMutex
	roomInfoReceivedLock   sync.RWMutex

	main    *roomViewMain
	toolbar *roomViewToolbar
	roster  *roomViewRoster
	conv    *roomViewConversation
	lobby   *roomViewLobby
}

func (v *roomView) onSelfJoinReceived(f func()) {
	v.selfJoinedReceivedLock.Lock()
	defer v.selfJoinedReceivedLock.Unlock()

	v.selfJoinedReceived = append(v.selfJoinedReceived, f)
}

func (v *roomView) onOccupantReceived(f func()) {
	v.occupantReceivedLock.Lock()
	defer v.occupantReceivedLock.Unlock()

	v.occupantReceived = append(v.occupantReceived, f)
}

func (v *roomView) onRoomInfoReceived(f func(*muc.RoomListing)) {
	v.roomInfoReceivedLock.Lock()
	defer v.roomInfoReceivedLock.Unlock()

	v.roomInfoReceived = append(v.roomInfoReceived, f)
}

func (a *account) newRoomView(ident jid.Bare, u *gtkUI) *roomView {
	view := &roomView{
		u:        u,
		account:  a,
		identity: ident,
		room:     a.newRoomModel(ident),
	}

	view.log = a.log.WithField("room", ident)

	view.initUIBuilder()
	view.initDefaults()

	toolbar := view.newRoomViewToolbar()
	view.toolbar = toolbar

	roster := view.newRoomViewRoster()
	view.roster = roster

	conversation := view.newRoomViewConversation()
	view.conv = conversation

	return view
}

func (v *roomView) setTitle(r string) {
	v.window.SetTitle(r)
}

func (u *gtkUI) getRoomOrCreateItIfNoExists(a *account, ident jid.Bare) (*roomView, bool) {
	v, ok := a.getRoomView(ident.String())
	if !ok {
		v = a.newRoomView(ident, u)
		a.addRoomView(v)
	}
	return v, ok
}

// mucShowRoom MUST be called always from the UI thread
//
// Also, when we want to show a chat room, having a "return to" function that
// will be called from the lobby only when the user wants to "cancel" or "return"
// might be useful in some scenarios like "returning to previous step".
//
// Please note that "returnTo" will be called from the UI thread too
func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare, returnTo func()) {
	view, ok := a.getRoomView(ident.String())
	if !ok {
		view = a.newRoomView(ident, u)
		a.addRoomView(view)
	}

	if view.joined {
		// In the main view of the room, we don't have the "cancel"
		// functionality that it's useful only in the lobby view of the room.
		// For that reason is why we ignore the "returnTo" value.
		view.returnTo = nil

		view.switchToMainView()
	} else {
		view.returnTo = returnTo
		view.switchToLobbyView()
	}

	if !ok {
		view.window.Show()
	} else {
		view.window.Present()
	}
	view.opened = true
}

func (v *roomView) isOpen() bool {
	return v.opened
}

func (v *roomView) initUIBuilder() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorNotif = newErrorNotification(v.notificationArea)

	v.builder.ConnectSignals(map[string]interface{}{
		"on_destroy_window": v.onDestroyWindow,
	})
}

func (v *roomView) initDefaults() {
	v.setTitle(v.identity.String())
	v.requestRoomInfo()
}

func (v *roomView) requestRoomInfo() {
	v.showSpinner()

	go func() {
		rl := make(chan *muc.RoomListing)
		go v.account.session.GetRoom(v.identity, rl)
		v.info = <-rl
		v.onRequestRoomInfoFinish()
	}()
}

func (v *roomView) onRequestRoomInfoFinish() {
	doInUIThread(v.hideSpinner)

	v.roomInfoReceivedLock.RLock()
	defer v.roomInfoReceivedLock.RUnlock()
	for _, f := range v.roomInfoReceived {
		f(v.info)
	}
}

func (v *roomView) onDestroyWindow() {
	v.opened = false
}

func (v *roomView) clearErrors() {
	v.errorNotif.Hide()
}

func (v *roomView) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}

	v.errorNotif.ShowMessage(err)
}

func (v *roomView) showSpinner() {
	v.spinner.Start()
	v.spinner.Show()
}

func (v *roomView) hideSpinner() {
	v.spinner.Stop()
	v.spinner.Hide()
}

func (v *roomView) tryLeaveRoom(onSuccess, onError func()) {
	v.clearErrors()
	v.showSpinner()

	go func() {
		v.account.leaveRoom(v.identity, v.occupant.Nick, func() {
			doInUIThread(v.window.Destroy)
			if onSuccess != nil {
				onSuccess()
			}
		}, func(err error) {
			//TODO: Should we use some notification manager?
			v.log.WithError(err).Error("An error occurred when trying to leave the room")
			doInUIThread(func() {
				v.hideSpinner()
				v.notifyOnError(i18n.Local("Couldn't leave the room, please try again."))
			})
			if onError != nil {
				onError()
			}
		})
	}()
}

func (v *roomView) switchToLobbyView() {
	if v.lobby == nil {
		v.lobby = v.newRoomViewLobby(v.account, v.identity, v.content, v.onJoined, v.onJoinCancel)
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
func (v *roomView) onRoomOccupantJoinedReceived(occupant string) {
	if v.joined {
		v.log.WithField("occupant", occupant).Error("A joined event was received but the user is already in the room")
		return
	}

	v.assignCurrentOccupant(v.identity.WithResource(jid.NewResource(occupant)).String())

	v.selfJoinedReceivedLock.RLock()
	defer v.selfJoinedReceivedLock.RUnlock()
	for _, f := range v.selfJoinedReceived {
		f()
	}
}

func (v *roomView) assignCurrentOccupant(occupantIdentity string) {
	o, ok := v.roster.r.GetOccupantByIdentity(occupantIdentity)
	if !ok {
		//TODO: Show in an appropriate way the error message to the user. Maybe with some ´handler notification´ struct?
		v.log.Error("An error occurred trying to get the current occupant")
		return
	}

	v.occupant = o
}

// onRoomOccupantUpdateReceived MUST be called from the UI thread
func (v *roomView) onRoomOccupantUpdateReceived() {
	v.occupantReceivedLock.RLock()
	defer v.occupantReceivedLock.RUnlock()
	for _, f := range v.occupantReceived {
		f()
	}
}

// onRoomOccupantLeftTheRoomReceived MUST be called from the UI thread
func (v *roomView) onRoomOccupantLeftTheRoomReceived(nickname string) {
	if v.conv != nil {
		v.conv.displayNotificationWhenOccupantLeftTheRoom(nickname)
	}
	v.roster.updateRosterModel()
}

// someoneJoinedTheRoom MUST be called from the UI thread
func (v *roomView) someoneJoinedTheRoom(nickname string) {
	if v.conv != nil {
		v.conv.displayNotificationWhenOccupantJoinedRoom(nickname)
	}
	v.roster.updateRosterModel()
}

// onRoomMessageToTheRoomReceived MUST be called from the UI thread
func (v *roomView) onRoomMessageToTheRoomReceived(nickname, subject, message string) {
	if v.conv != nil {
		v.conv.displayNewLiveMessage(nickname, subject, message)
	}
}

// loggingIsEnabled MUST not be called from the UI thread
func (v *roomView) loggingIsEnabled() {
	if v.conv != nil {
		msg := i18n.Local("This room is now publicly logged, meaning that everything you and the others in the room say or do can be made public on a website.")
		doInUIThread(func() {
			v.conv.displayWarningMessage(msg)
		})
	}
}

// loggingIsDisabled MUST not be called from the UI thread
func (v *roomView) loggingIsDisabled() {
	if v.conv != nil {
		msg := i18n.Local("This room is no longer publicly logged.")
		doInUIThread(func() {
			v.conv.displayWarningMessage(msg)
		})
	}
}
