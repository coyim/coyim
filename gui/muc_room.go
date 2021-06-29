package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewDataProvider interface {
	passwordProvider() string
	backToPreviousStep() func()
	roomProperties() *muc.RoomListing
}

type roomViewData struct {
	passsword            string
	onBackToPreviousStep func()
	properties           *muc.RoomListing
}

func newRoomViewData() *roomViewData {
	return &roomViewData{}
}

func (rvd *roomViewData) passwordProvider() string {
	return rvd.passsword
}

func (rvd *roomViewData) backToPreviousStep() func() {
	return rvd.onBackToPreviousStep
}

func (rvd *roomViewData) roomProperties() *muc.RoomListing {
	return rvd.properties
}

type roomView struct {
	u       *gtkUI
	account *account
	builder *builder

	room *muc.Room

	cancel chan bool

	opened             bool
	passwordProvider   func() string
	backToPreviousStep func()

	window                 gtki.Window   `gtk-widget:"room-window"`
	overlay                gtki.Overlay  `gtk-widget:"room-overlay"`
	privacityWarningBox    gtki.Box      `gtk-widget:"room-privacity-warnings-box"`
	loadingNotificationBox gtki.Box      `gtk-widget:"room-loading-notification-box"`
	content                gtki.Box      `gtk-widget:"room-main-box"`
	notificationsArea      gtki.Revealer `gtk-widget:"room-notifications-revealer"`
	roomInfoErrorBar       gtki.InfoBar  `gtk-widget:"room-info-error-bar"`

	notifications *roomNotifications

	warnings           *roomViewWarnings
	warningsInfoBar    *roomViewWarningsInfoBar
	loadingViewOverlay *roomViewLoadingOverlay

	subscribers *roomViewSubscribers

	main    *roomViewMain
	toolbar *roomViewToolbar
	roster  *roomViewRoster
	conv    *roomViewConversation
	lobby   *roomViewLobby

	log coylog.Logger
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
	view.initSubscribers()
	view.initNotifications()

	view.toolbar = view.newRoomViewToolbar()
	view.roster = view.newRoomViewRoster()
	view.conv = view.newRoomViewConversation()

	view.warnings = view.newRoomViewWarnings()
	view.warningsInfoBar = view.newRoomViewWarningsInfoBar()
	view.loadingViewOverlay = view.newRoomViewLoadingOverlay()

	view.requestRoomDiscoInfo()

	return view
}

func (v *roomView) initBuilderAndSignals() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_destroy_window":       v.onDestroyWindow,
		"on_room_info_load_retry": v.requestRoomDiscoInfo,
	})
}

func (v *roomView) initDefaults() {
	v.setTitle(i18n.Localf("%s [%s]", v.roomID(), v.account.Account()))

	mucStyles.setRoomWindowStyle(v.window)
}

func (v *roomView) initSubscribers() {
	v.subscribe("room", func(ev roomViewEvent) {
		doInUIThread(func() {
			v.onEventReceived(ev)
		})
	})
}

func (v *roomView) onEventReceived(ev roomViewEvent) {
	switch t := ev.(type) {
	case selfOccupantRemovedEvent:
		v.selfOccupantRemovedEvent()
	case roomDiscoInfoReceivedEvent:
		v.roomDiscoInfoReceivedEvent(t.info)
	case roomConfigRequestTimeoutEvent:
		v.roomConfigRequestTimeoutEvent()
	case selfOccupantAffiliationUpdatedEvent:
		v.selfOccupantAffiliationUpdatedEvent(t.selfAffiliationUpdate)
	case selfOccupantAffiliationRoleUpdatedEvent:
		v.selfOccupantAffiliationRoleUpdatedEvent(t.selfAffiliationRoleUpdate)
	case selfOccupantRoleUpdatedEvent:
		v.selfOccupantRoleUpdatedEvent(t.selfRoleUpdate)
	}
}

func (v *roomView) requestRoomDiscoInfo() {
	v.loadingViewOverlay.onRoomDiscoInfoLoad()
	v.roomInfoErrorBar.Hide()
	go v.account.session.GetRoomInformation(v.roomID())
}

// roomDiscoInfoReceivedEvent MUST be called from the UI thread
func (v *roomView) roomDiscoInfoReceivedEvent(di data.RoomDiscoInfo) {
	v.loadingViewOverlay.hide()

	v.warnings.clear()
	v.addRoomWarningsBasedOnInfo(di)
	v.privacityWarningBox.PackStart(v.warningsInfoBar.view(), true, false, 0)
}

// roomConfigRequestTimeoutEvent MUST be called from the UI thread
func (v *roomView) roomConfigRequestTimeoutEvent() {
	v.loadingViewOverlay.hide()
	v.warnings.clear()

	v.roomInfoErrorBar.Show()
}

// selfOccupantAffiliationUpdatedEvent MUST be called from the UI thread
func (v *roomView) selfOccupantAffiliationUpdatedEvent(selfAffiliationUpdate data.SelfAffiliationUpdate) {
	v.notifications.info(roomNotificationOptions{
		message:   getMUCNotificationMessageFrom(selfAffiliationUpdate),
		showTime:  true,
		closeable: true,
	})

	if selfAffiliationUpdate.New.IsBanned() {
		v.disableRoomView()
	}
}

// selfOccupantAffiliationRoleUpdatedEvent MUST be called from the UI thread
func (v *roomView) selfOccupantAffiliationRoleUpdatedEvent(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) {
	v.notifications.info(roomNotificationOptions{
		message:   getSelfAffiliationRoleUpdateMessage(selfAffiliationRoleUpdate),
		showTime:  true,
		closeable: true,
	})
}

// selfOccupantRoleUpdatedEvent MUST be called from the UI thread
func (v *roomView) selfOccupantRoleUpdatedEvent(selfRoleUpdate data.RoleUpdate) {
	v.notifications.info(roomNotificationOptions{
		message:   getSelfRoleUpdateMessage(selfRoleUpdate),
		showTime:  true,
		closeable: true,
	})

	if selfRoleUpdate.New.IsNone() {
		v.disableRoomView()
	}
}

// selfOccupantRemovedEvent MUST be called from the UI thread
func (v *roomView) selfOccupantRemovedEvent() {
	v.notifications.info(roomNotificationOptions{
		message:   i18n.Local("You have been removed from this room because it's now a members only room."),
		showTime:  true,
		closeable: true,
	})

	v.disableRoomView()
}

func (v *roomView) disableRoomView() {
	doInUIThread(func() {
		mucStyles.setDisableRoomStyle(v.content)
		v.account.removeRoomView(v.roomID())
		v.warningsInfoBar.hide()
	})
}

func (v *roomView) onDestroyWindow() {
	v.opened = false
	v.account.removeRoomView(v.roomID())
	go v.cancelActiveRequests()
}

// cancelActiveRequests MUST NOT be called from the UI thread
func (v *roomView) cancelActiveRequests() {
	if v.cancel != nil {
		v.cancel <- true
		v.cancel = nil
	}
}

func (v *roomView) setTitle(t string) {
	v.window.SetTitle(t)
}

func (v *roomView) isOpen() bool {
	return v.opened
}

func (v *roomView) isSelfOccupantInTheRoom() bool {
	return v.room.IsSelfOccupantInTheRoom()
}

func (v *roomView) isSelfOccupantAnOwner() bool {
	return v.room.IsSelfOccupantAnOwner()
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

func (v *roomView) onLeaveRoom() {
	// TODO: Implement the logic behind leaving this room and
	// how the view will interact with the user during this process
	v.tryLeaveRoom(nil, nil)
}

// tryLeaveRoom MUST be called from the UI thread.
// Please note that "onSuccess" and "onError" will be called from another thread.
func (v *roomView) tryLeaveRoom(onSuccess func(), onError func(error)) {
	onSuccessFinal := func() {
		doInUIThread(v.window.Destroy)
		if onSuccess != nil {
			onSuccess()
		}
	}

	onErrorFinal := func(err error) {
		v.log.WithError(err).Error("An error occurred when trying to leave the room")
		if onError != nil {
			onError(err)
		}
	}

	go v.account.leaveRoom(
		v.roomID(),
		v.room.SelfOccupantNickname(),
		onSuccessFinal,
		onErrorFinal,
		nil,
	)
}

func (v *roomView) publishDestroyEvent(reason string, alternativeRoomID jid.Bare, password string) {
	v.publishEvent(roomDestroyedEvent{
		reason:      reason,
		alternative: alternativeRoomID,
		password:    password,
	})
}

// tryDestroyRoom MUST be called from the UI thread, but please, note that
// the "onSuccess" and "onError" callbacks will be called from another thread
func (v *roomView) tryDestroyRoom(reason string, alternativeRoomID jid.Bare, password string) {
	v.loadingViewOverlay.onRoomDestroy()

	sc, ec := v.account.session.DestroyRoom(v.roomID(), reason, alternativeRoomID, password)
	go func() {
		select {
		case <-sc:
			v.log.Info("The room has been destroyed")
			v.publishDestroyEvent(reason, alternativeRoomID, password)
			doInUIThread(func() {
				v.loadingViewOverlay.hide()

				v.notifications.info(roomNotificationOptions{
					message:   i18n.Local("The room has been destroyed"),
					closeable: true,
				})
			})
		case err := <-ec:
			v.log.WithError(err).Error("An error occurred when trying to destroy the room")
			doInUIThread(func() {
				v.loadingViewOverlay.hide()

				dr := createDestroyDialogError(func() {
					v.tryDestroyRoom(reason, alternativeRoomID, password)
				})

				dr.updateErrorMessage(err)
				dr.show()
			})
		}
	}()
}

func (v *roomView) tryUpdateOccupantAffiliation(o *muc.Occupant, newAffiliation data.Affiliation, reason string) {
	doInUIThread(func() {
		v.loadingViewOverlay.onOccupantAffiliationUpdate()
	})

	previousAffiliation := o.Affiliation
	sc, ec := v.account.session.UpdateOccupantAffiliation(v.roomID(), o.Nickname, o.RealJid, newAffiliation, reason)

	select {
	case <-sc:
		v.log.Info("The affiliation has been changed")
		v.onOccupantAffiliationUpdateSuccess(o, previousAffiliation, newAffiliation)
	case err := <-ec:
		v.log.WithError(err).Error("An error occurred while updating the occupant affiliation")
		v.onOccupantAffiliationUpdateError(o.Nickname, newAffiliation, err)
	}
}

func (v *roomView) onOccupantAffiliationUpdateSuccess(o *muc.Occupant, previousAffiliation, affiliation data.Affiliation) {
	doInUIThread(func() {
		v.loadingViewOverlay.hide()

		v.notifications.info(roomNotificationOptions{
			message:   getAffiliationUpdateSuccessMessage(o.Nickname, previousAffiliation, affiliation),
			closeable: true,
		})
	})
}

func (v *roomView) onBannedListUpdated() {
	doInUIThread(func() {
		v.notifications.info(roomNotificationOptions{
			message:   i18n.Local("The banned list has been updated"),
			closeable: true,
		})
	})
}

func (v *roomView) onOccupantAffiliationUpdateError(nickname string, newAffiliation data.Affiliation, err error) {
	messages := getAffiliationUpdateFailureMessage(nickname, newAffiliation, err)

	doInUIThread(func() {
		v.loadingViewOverlay.hide()

		v.notifications.info(roomNotificationOptions{
			message:   messages.notificationMessage,
			closeable: true,
		})

		dr := createDialogErrorComponent(
			messages.errorDialogTitle,
			messages.errorDialogHeader,
			messages.errorDialogMessage,
		)

		dr.show()
	})
}

func (v *roomView) tryUpdateOccupantRole(o *muc.Occupant, newRole data.Role, reason string) {
	l := v.log.WithField("occupant", o.Nickname)

	doInUIThread(func() {
		v.loadingViewOverlay.onOccupantRoleUpdate(newRole)
	})

	previousRole := o.Role
	sc, ec := v.account.session.UpdateOccupantRole(v.roomID(), o.Nickname, newRole, reason)

	select {
	case <-sc:
		l.Info("The role has been changed")
		v.onOccupantRoleUpdateSuccess(o.Nickname, previousRole, newRole)
	case err := <-ec:
		l.WithError(err).Error("An error occurred while updating the occupant role")
		v.onOccupantRoleUpdateError(o.Nickname, newRole)
	}
}

func (v *roomView) onOccupantRoleUpdateSuccess(nickname string, previousRole, newRole data.Role) {
	doInUIThread(func() {
		v.loadingViewOverlay.hide()

		v.notifications.info(roomNotificationOptions{
			message:   getRoleUpdateSuccessMessage(nickname, previousRole, newRole),
			closeable: true,
		})
	})
}

func (v *roomView) onOccupantRoleUpdateError(nickname string, newRole data.Role) {
	messages := getRoleUpdateFailureMessage(nickname, newRole)

	doInUIThread(func() {
		v.loadingViewOverlay.hide()

		v.notifications.error(roomNotificationOptions{
			message:   messages.notificationMessage,
			closeable: true,
		})

		dr := createDialogErrorComponent(
			messages.errorDialogTitle,
			messages.errorDialogHeader,
			messages.errorDialogMessage,
		)

		dr.show()
	})
}

func (v *roomView) updateSubjectRoom(s string, onSuccess func()) {
	err := v.account.session.UpdateRoomSubject(v.roomID(), v.room.SelfOccupant().RealJid.String(), s)
	if err != nil {
		doInUIThread(func() {
			v.notifications.error(roomNotificationOptions{message: i18n.Local("The room subject couldn't be updated."), closeable: true})
		})
		return
	}
	doInUIThread(func() {
		onSuccess()

		v.notifications.info(roomNotificationOptions{
			message:   i18n.Local("The room subject has been updated."),
			closeable: true,
		})
	})
}

func (v *roomView) switchToLobbyView() {
	v.initRoomLobby()

	l := i18n.Local("Cancel")
	if v.backToPreviousStep != nil {
		l = i18n.Local("Return")
	}
	setFieldLabel(v.lobby.cancelButton, l)

	v.warningsInfoBar.onClose(nil)

	v.lobby.show()
}

func (v *roomView) switchToMainView() {
	v.initRoomMain()

	v.warningsInfoBar.onClose(v.warningsInfoBar.hide)

	v.main.show()
}

func (v *roomView) onJoined() {
	doInUIThread(func() {
		v.lobby.destroy()
		v.content.Remove(v.lobby.content)
		v.switchToMainView()
	})
}

func (v *roomView) onJoinCancel() {
	v.window.Destroy()

	if v.backToPreviousStep != nil {
		v.backToPreviousStep()
	}
}

// messageForbidden MUST NOT be called from the UI thread
func (v *roomView) messageForbidden() {
	v.publishEvent(messageForbidden{})
}

// messageNotAccepted MUST NOT be called from the UI thread
func (v *roomView) messageNotAccepted() {
	v.publishEvent(messageNotAcceptable{})
}

// nicknameConflict MUST NOT be called from the UI thread
func (v *roomView) nicknameConflict(nickname string) {
	v.publishEvent(nicknameConflictEvent{nickname})
}

// serviceUnavailableEvent MUST NOT be called from the UI thread
func (v *roomView) serviceUnavailable() {
	v.publishEvent(serviceUnavailableEvent{})
}

// unknownError MUST NOT be called from the UI thread
func (v *roomView) unknownError() {
	v.publishEvent(unknownErrorEvent{})
}

// registrationRequired MUST NOT be called from the UI thread
func (v *roomView) registrationRequired(nickname string) {
	v.publishEvent(registrationRequiredEvent{nickname})
}

// notAuthorized MUST NOT be called from the UI thread
func (v *roomView) notAuthorized() {
	v.publishEvent(notAuthorizedEvent{})
}

// occupantForbidden MUST NOT be called from the UI thread
func (v *roomView) occupantForbidden() {
	v.publishEvent(occupantForbiddenEvent{})
}

// publishOccupantAffiliationUpdatedEvent MUST NOT be called from the UI thread
func (v *roomView) publishOccupantAffiliationRoleUpdatedEvent(affiliationRoleUpdate data.AffiliationRoleUpdate) {
	v.publishEvent(occupantAffiliationRoleUpdatedEvent{affiliationRoleUpdate})
}

// publishSelfOccupantAffiliationUpdatedEvent MUST NOT be called from the UI thread
func (v *roomView) publishSelfOccupantAffiliationRoleUpdatedEvent(selfAffiliationRoleUpdate data.AffiliationRoleUpdate) {
	v.publishEvent(selfOccupantAffiliationRoleUpdatedEvent{selfAffiliationRoleUpdate})
}

// publishOccupantAffiliationUpdatedEvent MUST NOT be called from the UI thread
func (v *roomView) publishOccupantAffiliationUpdatedEvent(affiliationUpdate data.AffiliationUpdate) {
	v.publishEvent(occupantAffiliationUpdatedEvent{affiliationUpdate})
}

// publishSelfOccupantAffiliationUpdatedEvent MUST NOT be called from the UI thread
func (v *roomView) publishSelfOccupantAffiliationUpdatedEvent(selfAffiliationUpdate data.SelfAffiliationUpdate) {
	v.publishEvent(selfOccupantAffiliationUpdatedEvent{selfAffiliationUpdate})
}

// publishOccupantRoleUpdatedEvent MUST NOT be called from the UI thread
func (v *roomView) publishOccupantRoleUpdatedEvent(roleUpdate data.RoleUpdate) {
	v.publishEvent(occupantRoleUpdatedEvent{roleUpdate})
}

// publishSelfOccupantRoleUpdatedEvent MUST NOT be called from the UI thread
func (v *roomView) publishSelfOccupantRoleUpdatedEvent(selfRoleUpdate data.RoleUpdate) {
	v.publishEvent(selfOccupantRoleUpdatedEvent{selfRoleUpdate})
}

// mainWindow MUST be called from the UI thread
func (v *roomView) mainWindow() gtki.Window {
	return v.window
}

func (v *roomView) roomID() jid.Bare {
	return v.room.ID
}

func (v *roomView) roomDisplayName() string {
	return v.roomID().Local().String()
}
