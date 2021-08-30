package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

type roomViewLobby struct {
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

	errorNotifications canNotifyErrors

	log coylog.Logger
}

func (v *roomView) newRoomViewLobby() *roomViewLobby {
	l := &roomViewLobby{
		roomView:              v,
		account:               v.account,
		errorNotifications:    v.notifications,
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
	l.roomNameLabel.SetText(i18n.Localf("You are joining %s", l.roomView.roomID()))
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
		case roomDisableEvent:
			l.disableLobbyFields()
		}
	})
}

func (l *roomViewLobby) roomDiscoInfoReceivedEvent(di data.RoomDiscoInfo, passwordProvider func() string) {
	l.isReadyToJoinRoom = true
	doInUIThread(func() {
		l.isPasswordProtected = di.PasswordProtected

		l.enableLobbyFields()
		l.enableJoinIfConditionsAreMet()

		if l.isPasswordProtected && passwordProvider != nil {
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

// isNotNicknameInConflictList MUST be called from the UI thread
func (l *roomViewLobby) isNotNicknameInConflictList() bool {
	if l.nicknamesWithConflict.Has(getEntryText(l.nicknameEntry)) {
		l.errorNotifications.notifyOnError(i18n.Local("That nickname is already being used."))
		return false
	}
	return true
}

// enableJoinIfConditionsAreMet MUST be called from the UI thread
func (l *roomViewLobby) enableJoinIfConditionsAreMet() {
	if !l.accountIsBanned {
		l.errorNotifications.clearErrors()
		setFieldSensitive(l.joinButton, l.checkJoinConditions())
	}
}

// checkJoinConditions MUST be called from the UI thread
func (l *roomViewLobby) checkJoinConditions() bool {
	nicknameHasContent := getEntryText(l.nicknameEntry) != ""
	passwordHasContent := getEntryText(l.passwordEntry) != ""

	return l.isReadyToJoinRoom && nicknameHasContent && l.isNotNicknameInConflictList() &&
		(!l.isPasswordProtected || passwordHasContent)
}

// disableLobbyFields MUST be called from the UI thread
func (l *roomViewLobby) disableLobbyFields() {
	disableField(l.nicknameEntry)
	disableField(l.passwordEntry)
	disableField(l.joinButton)
}

// enableLobbyFields MUST be called from the UI thread
func (l *roomViewLobby) enableLobbyFields() {
	enableField(l.nicknameEntry)
	setFieldSensitive(l.passwordEntry, l.isPasswordProtected)
	enableField(l.joinButton)

	setFieldVisibility(l.passwordLabel, l.isPasswordProtected)
	setFieldVisibility(l.passwordEntry, l.isPasswordProtected)
}

// onJoinRoomClicked MUST be called from the UI thread
func (l *roomViewLobby) onJoinRoomClicked(done func()) {
	l.errorNotifications.clearErrors()
	l.disableLobbyFields()

	nickname := getEntryText(l.nicknameEntry)
	password := getEntryText(l.passwordEntry)

	l.roomView.loadingViewOverlay.onJoinRoom()
	go l.roomView.sendJoinRoomRequest(nickname, password, done)
}

// onJoinFailed MUST be called from the UI thread
func (l *roomViewLobby) onJoinFailed(err error) {
	l.enableLobbyFields()
	shouldEnableJoin := l.checkJoinConditions()

	userMessage := i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
	if err, ok := err.(*roomError); ok {
		userMessage = err.friendlyMessage
		shouldEnableJoin = false

		if err.errType == errJoinNicknameConflict {
			if !l.nicknamesWithConflict.Has(err.nickname) {
				l.nicknamesWithConflict.Insert(err.nickname)
			}
		}
	}

	l.errorNotifications.notifyOnError(userMessage)

	setFieldSensitive(l.joinButton, shouldEnableJoin)
}
