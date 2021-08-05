package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

type roomViewLobby struct {
	roomID                jid.Bare
	roomView              *roomView
	account               *account
	accountIsBanned       bool
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

	notifications  *roomNotifications
	loadingOverlay *roomViewLoadingOverlay

	log coylog.Logger
}

func (v *roomView) initRoomLobby() {
	v.lobby = v.newRoomViewLobby(v.account, v.roomID())
	v.content.Add(v.lobby.content)
}

func (v *roomView) newRoomViewLobby(a *account, roomID jid.Bare) *roomViewLobby {
	l := &roomViewLobby{
		roomID:                roomID,
		roomView:              v,
		account:               a,
		loadingOverlay:        v.loadingViewOverlay,
		nicknamesWithConflict: set.New(),
		log:                   v.log.WithField("where", "roomViewLobby"),
	}

	l.initBuilder()
	l.initDefaults()
	l.initSubscribers()

	return l
}

func (l *roomViewLobby) initBuilder() {
	builder := newBuilder("MUCRoomLobby")
	panicOnDevError(builder.bindObjects(l))

	builder.ConnectSignals(map[string]interface{}{
		"on_nickname_changed": l.enableJoinIfConditionsAreMet,
		"on_password_changed": l.enableJoinIfConditionsAreMet,
		"on_join":             doOnlyOnceAtATime(l.onJoinRoomClicked),
		"on_cancel":           l.roomView.onJoinCancel,
	})
}

func (l *roomViewLobby) initDefaults() {
	l.notifications = l.roomView.notifications

	l.roomNameLabel.SetText(i18n.Localf("You are joining %s", l.roomID.String()))
	l.content.SetHExpand(true)

	setFieldVisibility(l.passwordLabel, false)
	setFieldVisibility(l.passwordEntry, false)

	mucStyles.setRoomToolbarLobyStyle(l.content)
}

func (l *roomViewLobby) initSubscribers() {
	l.roomView.subscribe("lobby", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case roomDiscoInfoReceivedEvent:
			l.roomDiscoInfoReceivedEvent(t.info, l.roomView.passwordProvider)
		case roomConfigRequestTimeoutEvent:
			l.roomConfigRequestTimeoutEvent()
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

func (l *roomViewLobby) roomConfigRequestTimeoutEvent() {
	l.isReadyToJoinRoom = false
	doInUIThread(func() {
		disableField(l.nicknameEntry)
		disableField(l.passwordEntry)
	})
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
		l.notifications.error(roomNotificationOptions{message: i18n.Local("That nickname is already being used.")})
		return false
	}
	return true
}

func (l *roomViewLobby) enableJoinIfConditionsAreMet() {
	if !l.accountIsBanned {
		l.notifications.clearAll()
		setFieldSensitive(l.joinButton, l.checkJoinConditions())
	}
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

// onJoinRoomClicked MUST be called from the UI thread
func (l *roomViewLobby) onJoinRoomClicked(done func()) {
	l.notifications.clearAll()
	l.disableFieldsAndShowSpinner()

	nickname := getEntryText(l.nicknameEntry)
	password := getEntryText(l.passwordEntry)

	go l.roomView.sendJoinRoomRequest(nickname, password, done)
}

func (l *roomViewLobby) onJoinFailed(err error) {
	l.enableFieldsAndHideSpinner()
	shouldEnableJoin := l.checkJoinConditions()

	userMessage := i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
	if err, ok := err.(*roomError); ok {
		userMessage = l.getUserErrorMessage(err)
		shouldEnableJoin = false

		if err.errType == errJoinNicknameConflict {
			if !l.nicknamesWithConflict.Has(err.nickname) {
				l.nicknamesWithConflict.Insert(err.nickname)
			}
		}
	}

	l.notifications.error(roomNotificationOptions{message: userMessage})

	setFieldSensitive(l.joinButton, shouldEnableJoin)
}
