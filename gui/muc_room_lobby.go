package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

type roomViewLobby struct {
	roomID                jid.Bare
	account               *account
	isPasswordProtected   bool
	isReadyToJoinRoom     bool
	nicknamesWithConflict *set.Set

	content          gtki.Box    `gtk-widget:"main-content"`
	roomNameLabel    gtki.Label  `gtk-widget:"room-name-value"`
	nicknameEntry    gtki.Entry  `gtk-widget:"nickname-entry"`
	passwordLabel    gtki.Label  `gtk-widget:"password-label"`
	passwordEntry    gtki.Entry  `gtk-widget:"password-entry"`
	joinButton       gtki.Button `gtk-widget:"join-button"`
	cancelButton     gtki.Button `gtk-widget:"cancel-button"`
	notificationArea gtki.Box    `gtk-widget:"notifications-box"`

	notifications  *notificationsComponent
	loadingOverlay *roomViewLoadingOverlay

	// These two methods WILL BE called from the UI thread
	onSuccess func()
	onCancel  func()

	log coylog.Logger
}

func (v *roomView) initRoomLobby() {
	v.lobby = v.newRoomViewLobby(v.account, v.roomID())
	v.content.Add(v.lobby.content)
}

func (v *roomView) newRoomViewLobby(a *account, roomID jid.Bare) *roomViewLobby {
	l := &roomViewLobby{
		roomID:                roomID,
		account:               a,
		onSuccess:             v.onJoined,
		onCancel:              v.onJoinCancel,
		loadingOverlay:        v.loadingViewOverlay,
		nicknamesWithConflict: set.New(),
		log:                   v.log.WithField("where", "roomViewLobby"),
	}

	l.initBuilder()
	l.initDefaults(v)
	l.initSubscribers(v)

	return l
}

func (l *roomViewLobby) initBuilder() {
	builder := newBuilder("MUCRoomLobby")
	panicOnDevError(builder.bindObjects(l))

	builder.ConnectSignals(map[string]interface{}{
		"on_nickname_changed": l.enableJoinIfConditionsAreMet,
		"on_password_changed": l.enableJoinIfConditionsAreMet,
		"on_join":             doOnlyOnceAtATime(l.onJoinRoom),
		"on_cancel":           l.onCancel,
	})
}

func (l *roomViewLobby) initDefaults(v *roomView) {
	l.notifications = v.u.newNotificationsComponent()
	l.notificationArea.Add(l.notifications.getBox())

	l.roomNameLabel.SetText(l.roomID.String())
	l.content.SetHExpand(true)

	setFieldVisibility(l.passwordLabel, false)
	setFieldVisibility(l.passwordEntry, false)
}

func (l *roomViewLobby) initSubscribers(v *roomView) {
	v.subscribe("lobby", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case roomDiscoInfoReceivedEvent:
			l.roomDiscoInfoReceivedEvent(t.info, v.passwordProvider)
		case occupantSelfJoinedEvent:
			l.finishJoinRequest()
		case nicknameConflictEvent:
			l.nicknameConflictEvent(t.nickname)
		case registrationRequiredEvent:
			l.registrationRequiredEvent()
		case notAuthorizedEvent:
			l.notAuthorizedEvent()
		case serviceUnavailableEvent:
			l.serviceUnavailableEvent()
		case unknownErrorEvent:
			l.unknownErrorEvent()
		case occupantForbiddenEvent:
			l.occupantForbiddenEvent()
		}
	})
}

func (l *roomViewLobby) roomDiscoInfoReceivedEvent(di data.RoomDiscoInfo, passwordProvider func() string) {
	l.isReadyToJoinRoom = true
	doInUIThread(func() {
		l.enableJoinIfConditionsAreMet()
		if di.PasswordProtected {
			l.isPasswordProtected = true
			setFieldVisibility(l.passwordLabel, true)
			setFieldVisibility(l.passwordEntry, true)
			setEntryText(l.passwordEntry, passwordProvider())
		}
	})
}

func (l *roomViewLobby) finishJoinRequest() {
	doInUIThread(func() {
		l.notifications.clearAll()
		l.onSuccess()
	})
}

func (l *roomViewLobby) finishJoinRequestWithError(err error) {
	l.log.WithError(err).Error("An error occurred while trying to join the room")
	doInUIThread(func() {
		l.onJoinFailed(err)
	})
}

func (l *roomViewLobby) switchToReturnOnCancel() {
	setFieldLabel(l.cancelButton, i18n.Local("Return"))
}

func (l *roomViewLobby) switchToCancel() {
	setFieldLabel(l.cancelButton, i18n.Local("Cancel"))
}

func (l *roomViewLobby) show() {
	l.content.Show()
}

func (l *roomViewLobby) destroy() {
	l.loadingOverlay.hide()
	l.content.Destroy()
}

func (l *roomViewLobby) nicknameHasContent() bool {
	return getEntryText(l.nicknameEntry) != ""
}

func (l *roomViewLobby) passwordHasContent() bool {
	return getEntryText(l.passwordEntry) != ""
}

func (l *roomViewLobby) isNotNicknameInConflictList() bool {
	if l.nicknamesWithConflict.Has(getEntryText(l.nicknameEntry)) {
		l.notifications.error(i18n.Local("That nickname is already being used."))
		return false
	}
	return true
}

func (l *roomViewLobby) enableJoinIfConditionsAreMet() {
	l.notifications.clearErrors()
	setFieldSensitive(l.joinButton, l.checkJoinConditions())
}

func (l *roomViewLobby) checkJoinConditions() bool {
	return l.isReadyToJoinRoom && l.nicknameHasContent() && l.isNotNicknameInConflictList() &&
		(!l.isPasswordProtected || l.passwordHasContent())
}

func (l *roomViewLobby) disableFieldsAndShowSpinner() {
	disableField(l.nicknameEntry)
	disableField(l.joinButton)
	l.loadingOverlay.onJoinRoom()
}

func (l *roomViewLobby) enableFieldsAndHideSpinner() {
	enableField(l.nicknameEntry)
	enableField(l.joinButton)
	l.loadingOverlay.hide()
}

func (l *roomViewLobby) onJoinRoom(done func()) {
	nickname := getEntryText(l.nicknameEntry)
	password := getEntryText(l.passwordEntry)

	go l.sendJoinRoomRequest(nickname, password, done)
}

func (l *roomViewLobby) sendJoinRoomRequest(nickname, password string, done func()) {
	defer done()

	err := l.joinRoom(nickname, password)
	if err != nil {
		l.finishJoinRequestWithError(err)
		return
	}

	doInUIThread(l.disableFieldsAndShowSpinner)
}

func (l *roomViewLobby) joinRoom(nickname, password string) error {
	err := l.account.session.JoinRoom(l.roomID, nickname, password)
	if err == session.ErrMUCJoinRoomInvalidNickname {
		return newRoomLobbyInvalidNicknameError()
	}
	return err
}

func (l *roomViewLobby) onJoinFailed(err error) {
	l.enableFieldsAndHideSpinner()
	shouldEnableJoin := l.checkJoinConditions()

	userMessage := i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
	if err, ok := err.(*mucRoomLobbyErr); ok {
		userMessage = l.getUserErrorMessage(err)

		if err.errType == errJoinNicknameConflict {
			shouldEnableJoin = false
			if !l.nicknamesWithConflict.Has(err.nickname) {
				l.nicknamesWithConflict.Insert(err.nickname)
			}
		}
	}

	l.notifications.error(userMessage)

	setFieldSensitive(l.joinButton, shouldEnableJoin)
}
