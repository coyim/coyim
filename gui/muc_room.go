package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomView struct {
	u       *gtkUI
	account *account
	builder *builder

	identity jid.Bare

	room *muc.Room
	info *muc.RoomListing

	log      coylog.Logger
	opened   bool
	returnTo func()

	events chan muc.MUC

	window           gtki.Window  `gtk-widget:"roomWindow"`
	content          gtki.Box     `gtk-widget:"boxMainView"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"roomNotificationArea"`

	notification gtki.InfoBar
	errorNotif   *errorNotification

	subscribers *roomViewSubscribers

	main    *roomViewMain
	toolbar *roomViewToolbar
	roster  *roomViewRoster
	conv    *roomViewConversation
	lobby   *roomViewLobby

	sync.Mutex
}

func newRoomView(u *gtkUI, a *account, ident jid.Bare) *roomView {
	view := &roomView{
		u:        u,
		account:  a,
		identity: ident,
		events:   make(chan muc.MUC),
	}

	view.room = a.newRoomModel(ident)
	view.log = a.log.WithField("room", ident)

	view.room.Subscribe(view.events)
	go view.observeRoomEvents()

	view.subscribers = newRoomViewSubscribers(view.identity, view.log)

	view.initBuilderAndSignals()
	view.initDefaults()

	view.toolbar = view.newRoomViewToolbar()
	view.roster = view.newRoomViewRoster()
	view.conv = view.newRoomViewConversation()

	go view.requestRoomInfo()

	return view
}

func (v *roomView) initBuilderAndSignals() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.errorNotif = newErrorNotification(v.notificationArea)

	v.builder.ConnectSignals(map[string]interface{}{
		"on_destroy_window": v.onDestroyWindow,
	})
}

func (v *roomView) initDefaults() {
	v.setTitle(fmt.Sprintf("%s [%s]", v.identity.String(), v.account.Account()))
}

func (v *roomView) requestRoomInfo() {
	doInUIThread(v.showSpinner)

	rl := make(chan *muc.RoomListing)
	go v.account.session.GetRoom(v.identity, rl)
	go func() {
		roomInfo := <-rl
		v.onRequestRoomInfoFinish(roomInfo)
	}()
}

func (v *roomView) onRequestRoomInfoFinish(roomInfo *muc.RoomListing) {
	v.Lock()
	defer v.Unlock()

	v.info = roomInfo

	doInUIThread(func() {
		v.hideSpinner()
		v.publish(roomInfoReceived)
	})
}

func (v *roomView) onDestroyWindow() {
	v.opened = false
	v.account.removeRoomView(v.identity)
}

func (v *roomView) setTitle(t string) {
	v.window.SetTitle(t)
}

func (v *roomView) isOpen() bool {
	return v.opened
}

func (v *roomView) isJoined() bool {
	return v.room.Joined
}

func (v *roomView) present() {
	if v.isOpen() {
		v.window.Present()
	}
}

func (v *roomView) show() {
	if v.isOpen() {
		v.log.Debug("show(): the room view is already opened")
		return
	}

	v.opened = true
	v.window.Show()
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
	if !v.room.Joined {
		v.log.Debug("tryLeaveRoom(): trying to leave a not joined room")
		doInUIThread(func() {
			v.notifyOnError(i18n.Local("Couldn't leave the room, please try again."))
		})
		return
	}

	v.clearErrors()
	v.showSpinner()

	go func() {
		v.account.leaveRoom(v.identity, v.room.Occupant.Nick, func() {
			doInUIThread(v.window.Destroy)
			if onSuccess != nil {
				onSuccess()
			}
		}, func(err error) {
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
	v.publish(previousToSwitchToLobby)

	v.initRoomLobby()

	if v.shouldReturnOnCancel() {
		v.lobby.swtichToReturnOnCancel()
	} else {
		v.lobby.swtichToCancel()
	}

	v.lobby.show()
}

func (v *roomView) switchToMainView() {
	v.publish(previousToSwitchToMain)
	v.initRoomMain()
	v.main.show()
}

func (v *roomView) onJoined() {
	doInUIThread(func() {
		v.lobby.hide()
		v.switchToMainView()
	})
}

func (v *roomView) shouldReturnOnCancel() bool {
	return v.returnTo != nil
}

// TODO: if we have an active connection or request, we should
// stop/close it here before destroying the window
func (v *roomView) onJoinCancel() {
	v.window.Destroy()

	if v.shouldReturnOnCancel() {
		v.returnTo()
	}
}

// nicknameConflict MUST not be called from the UI thread
func (v *roomView) nicknameConflict(nickname string) {
	v.publishWithInfo(nicknameConflict, roomViewEventInfo{
		nickname: nickname,
	})
}

// registrationRequired MUST not be called from the UI thread
func (v *roomView) registrationRequired(nickname string) {
	v.publishWithInfo(registrationRequired, roomViewEventInfo{
		nickname: nickname,
	})
}
