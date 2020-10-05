package gui

import (
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

	room     *muc.Room
	info     *muc.RoomListing
	infoLock sync.Mutex

	log      coylog.Logger
	opened   bool
	returnTo func()

	window                 gtki.Window  `gtk-widget:"roomWindow"`
	content                gtki.Box     `gtk-widget:"boxMainView"`
	messagesOverlay        gtki.Overlay `gtk-widget:"messagesOverlay"`
	messagesOverlayBox     gtki.Box     `gtk-widget:"messagesOverlayBox"`
	messagesBox            gtki.Box     `gtk-widget:"messagesBox"`
	notificationBox        gtki.Box     `gtk-widget:"notificationBox"`
	loadingNotificationBox gtki.Box     `gtk-widget:"loadingNotificationBox"`

	spinner         *spinner
	notifications   *notifications
	loadingInfoBar  *roomViewLoadingInfoBar
	warnings        *roomViewWarningsOverlay
	warningsInfoBar *roomViewWarningsInfoBar

	subscribers *roomViewSubscribers

	main    *roomViewMain
	toolbar *roomViewToolbar
	roster  *roomViewRoster
	conv    *roomViewConversation
	lobby   *roomViewLobby
}

func newRoomView(u *gtkUI, a *account, roomID jid.Bare) *roomView {
	view := &roomView{
		u:       u,
		account: a,
	}

	// TODO: We already know this need to change
	view.room = a.newRoomModel(roomID)
	view.log = a.log.WithField("room", roomID)

	view.room.Subscribe(view.handleRoomEvent)

	view.subscribers = newRoomViewSubscribers(view.roomID(), view.log)

	view.initBuilderAndSignals()
	view.initDefaults()

	view.toolbar = view.newRoomViewToolbar()
	view.roster = view.newRoomViewRoster()
	view.conv = view.newRoomViewConversation()

	view.spinner = newSpinner()
	view.notifications = u.newNotifications(view.notificationBox)
	view.loadingInfoBar = view.newRoomViewLoadingInfoBar(view.loadingNotificationBox)
	view.warnings = view.newRoomViewWarningsOverlay(view.closeNotificationsOverlay)
	view.warningsInfoBar = view.newRoomViewWarningsInfoBar(view.showWarnings, view.removeWarningsInfobar)

	go view.requestRoomInfo()

	return view
}

func (v *roomView) initBuilderAndSignals() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_destroy_window": v.onDestroyWindow,
	})
}

func (v *roomView) initDefaults() {
	v.setTitle(i18n.Localf("%s [%s]", v.roomID(), v.account.Account()))
}

func (v *roomView) requestRoomInfo() {
	doInUIThread(v.loadingInfoBar.show)

	rl := make(chan *muc.RoomListing)
	go v.account.session.GetRoom(v.roomID(), rl)
	go func() {
		// TODO: What happens if no result ever comes? Maybe we need a timeout
		roomInfo := <-rl
		v.onRequestRoomInfoFinish(roomInfo)
	}()
}

func (v *roomView) onRequestRoomInfoFinish(roomInfo *muc.RoomListing) {
	v.infoLock.Lock()
	defer v.infoLock.Unlock()

	v.info = roomInfo

	doInUIThread(func() {
		v.warnings.clear()
		v.showRoomWarnings()
		v.notifications.add(v.warningsInfoBar)
	})

	doInUIThread(v.loadingInfoBar.hide)

	v.publish("roomInfoReceivedEvent")
}

func (v *roomView) showRoomWarnings() {
	v.warnings.add(i18n.Local("Please be aware that communication in chat rooms is not encrypted - anyone that can intercept communication between you and the server - and the server itself - will be able to see what you are saying in this chat room. Only join this room and communicate here if you trust the server to not be hostile."))

	switch v.info.Anonymity {
	case "semi":
		v.warnings.add(i18n.Local("This room is partially anonymous. This means that only moderators can connect your nickname with your real username (your JID)."))
	case "no":
		v.warnings.add(i18n.Local("This room is not anonymous. This means that any person in this room can connect your nickname with your real username (your JID)."))
	default:
		v.log.WithField("anonymity", v.info.Anonymity).Warn("Unknown anonymity setting for room")
	}

	if v.info.Logged {
		v.warnings.add(i18n.Local("This room is publicly logged, meaning that everything you and the others in the room say or do can be made public on a website."))
	}
}

func (v *roomView) showWarnings() {
	prov := providerWithStyle("box", style{
		"background-color": "#ffffff",
		"box-shadow":       "0 10px 20px rgba(0, 0, 0, 0.35)",
	})

	updateWithStyle(v.messagesBox, prov)

	v.warnings.show()
	v.showNotificationsOverlay()
}

func (v *roomView) removeWarningsInfobar() {
	v.notifications.remove(v.warningsInfoBar.getWidget())
}

func (v *roomView) showNotificationsOverlay() {
	prov := providerWithStyle("box", style{
		"background-color": "rgba(0, 0, 0, 0.5)",
	})

	updateWithStyle(v.messagesOverlayBox, prov)

	v.messagesOverlay.Show()
}

func (v *roomView) closeNotificationsOverlay() {
	v.messagesOverlay.Hide()
}

func (v *roomView) onDestroyWindow() {
	v.opened = false
	v.account.removeRoomView(v.roomID())
}

func (v *roomView) setTitle(t string) {
	v.window.SetTitle(t)
}

func (v *roomView) isOpen() bool {
	return v.opened
}

func (v *roomView) isJoined() bool {
	return v.room.IsSelfOccupantJoined()
}

func (v *roomView) present() {
	if v.isOpen() {
		v.window.Present()
	}
}

func (v *roomView) show() {
	v.opened = true
	v.window.Show()
}

func (v *roomView) tryLeaveRoom(onSuccess, onError func()) {
	v.spinner.show()

	go func() {
		v.account.leaveRoom(v.roomID(), v.room.SelfOccupantNickname(), func() {
			doInUIThread(v.window.Destroy)
			if onSuccess != nil {
				onSuccess()
			}
		}, func(err error) {
			v.log.WithError(err).Error("An error occurred when trying to leave the room")
			doInUIThread(v.spinner.hide)
			if onError != nil {
				onError()
			}
		})
	}()
}

func (v *roomView) switchToLobbyView() {
	v.publish("beforeSwitchingToLobbyViewEvent")

	v.initRoomLobby()

	if v.shouldReturnOnCancel() {
		v.lobby.switchToReturnOnCancel()
	} else {
		v.lobby.switchToCancel()
	}

	v.lobby.show()
}

func (v *roomView) switchToMainView() {
	v.publish("beforeSwitchingToMainViewEvent")
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

// nicknameConflict MUST NOT be called from the UI thread
func (v *roomView) nicknameConflict(nickname string) {
	v.publishWithInfo("nicknameConflictEvent", roomViewEventInfo{
		"nickname": nickname,
	})
}

// registrationRequired MUST NOT be called from the UI thread
func (v *roomView) registrationRequired(nickname string) {
	v.publishWithInfo("registrationRequiredEvent", roomViewEventInfo{
		"nickname": nickname,
	})
}

func (v *roomView) roomID() jid.Bare {
	return v.room.ID
}
