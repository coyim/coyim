package gui

import (
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc/data"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/golang-collections/collections/set"
)

type roomViewLobby struct {
	roomID  jid.Bare
	account *account
	log     coylog.Logger

	content          gtki.Box     `gtk-widget:"boxJoinRoomView"`
	mainBox          gtki.Box     `gtk-widget:"mainContent"`
	roomNameLabel    gtki.Label   `gtk-widget:"roomNameValue"`
	nicknameEntry    gtki.Entry   `gtk-widget:"nicknameEntry"`
	passwordLabel    gtki.Label   `gtk-widget:"passwordLabel"`
	passwordEntry    gtki.Entry   `gtk-widget:"passwordEntry"`
	joinButton       gtki.Button  `gtk-widget:"joinButton"`
	cancelButton     gtki.Button  `gtk-widget:"cancelButton"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"boxNotificationArea"`
	warningsArea     gtki.Box     `gtk-widget:"boxWarningsArea"`
	parent           gtki.Box
	errorNotif       *errorNotification
	notification     gtki.InfoBar

	onJoin                  chan bool
	onJoinError             chan error
	cancel                  chan bool
	roomIsPasswordProtected bool

	isReadyToJoinRoom bool

	nicknamesWithConflict *set.Set

	// These two methods WILL BE called from the UI thread
	onSuccess func()
	onCancel  func()
}

func (v *roomView) initRoomLobby() {
	v.lobby = v.newRoomViewLobby(v.account, v.roomID(), v.content, v.onJoined, v.onJoinCancel)
}

func (v *roomView) newRoomViewLobby(a *account, roomID jid.Bare, parent gtki.Box, onSuccess, onCancel func()) *roomViewLobby {
	l := &roomViewLobby{
		roomID:                roomID,
		account:               a,
		parent:                parent,
		onCancel:              onCancel,
		nicknamesWithConflict: set.New(),
		log:                   v.log.WithField("who", "roomViewLobby"),
	}

	l.onSuccess = func() {
		if onSuccess != nil {
			onSuccess()
		}
	}

	l.onCancel = func() {
		if onCancel != nil {
			onCancel()
		}
	}

	l.initBuilder()
	l.initDefaults()
	l.initSubscribers(v)

	return l
}

func (l *roomViewLobby) initBuilder() {
	builder := newBuilder("MUCRoomLobby")
	panicOnDevError(builder.bindObjects(l))

	builder.ConnectSignals(map[string]interface{}{
		"on_nickname_changed": l.onNicknameChange,
		"on_password_changed": l.onPasswordChange,
		"on_joined_clicked":   l.onJoinClicked,
		"on_cancel_clicked":   l.onJoinCancel,
	})
}

func (l *roomViewLobby) initDefaults() {
	l.errorNotif = newErrorNotification(l.notificationArea)

	l.roomNameLabel.SetText(l.roomID.String())
	l.content.SetHExpand(true)
	l.parent.Add(l.content)
	l.content.SetCenterWidget(l.mainBox)
}

func (l *roomViewLobby) initSubscribers(v *roomView) {
	v.subscribe("lobby", func(ev roomViewEvent) {
		switch t := ev.(type) {
		case occupantSelfJoinedEvent:
			l.occupantSelfJoinedEvent()
		case roomDiscoInfoReceivedEvent:
			l.roomDiscoInfoReceivedEvent(t.info)
		case nicknameConflictEvent:
			l.nicknameConflictEvent(l.roomID, t.nickname)
		case registrationRequiredEvent:
			l.registrationRequiredEvent(l.roomID, t.nickname)
		case notAuthorizedEvent:
			l.notAuthorizedEvent()
		}
	})
}

func (l *roomViewLobby) roomDiscoInfoReceivedEvent(di data.RoomDiscoInfo) {
	l.isReadyToJoinRoom = true
	doInUIThread(func() {
		l.enableJoinIfConditionsAreMet()
		if di.PasswordProtected {
			l.roomIsPasswordProtected = true
			l.passwordLabel.SetVisible(true)
			l.passwordEntry.SetVisible(true)
		}
	})
}

func (l *roomViewLobby) switchToReturnOnCancel() {
	l.cancelButton.SetProperty("label", i18n.Local("Return"))
}

func (l *roomViewLobby) switchToCancel() {
	l.cancelButton.SetProperty("label", i18n.Local("Cancel"))
}

func (l *roomViewLobby) show() {
	l.content.Show()
}

func (l *roomViewLobby) hide() {
	l.content.Hide()
}

func (l *roomViewLobby) close() {
	l.hide()
	l.parent.Remove(l.content)
}

func (l *roomViewLobby) onNicknameChange() {
	l.enableJoinIfConditionsAreMet()
}

func (l *roomViewLobby) onPasswordChange() {
	l.enableJoinIfConditionsAreMet()
}

func (l *roomViewLobby) nicknameHasContent() bool {
	nickname, _ := l.nicknameEntry.GetText()
	return nickname != ""
}

func (l *roomViewLobby) passwordHasContent() bool {
	password, _ := l.passwordEntry.GetText()
	return password != ""
}

func (l *roomViewLobby) isNotNicknameInConflictList() bool {
	nickname, _ := l.nicknameEntry.GetText()
	if l.nicknamesWithConflict.Has(nickname) {
		l.notifyOnError(i18n.Local("That nickname is already being used."))
		return false
	}
	return true
}

func (l *roomViewLobby) enableJoinIfConditionsAreMet() {
	l.clearErrors()

	conditionsAreValid := l.isReadyToJoinRoom && l.nicknameHasContent() && l.isNotNicknameInConflictList()
	if l.roomIsPasswordProtected {
		conditionsAreValid = conditionsAreValid && l.passwordHasContent()
	}

	l.joinButton.SetSensitive(conditionsAreValid)
}

func (l *roomViewLobby) disableFields() {
	l.nicknameEntry.SetSensitive(false)
}

func (l *roomViewLobby) enableFields() {
	l.nicknameEntry.SetSensitive(true)
}

func (l *roomViewLobby) showSpinner() {
	l.spinner.Start()
	l.spinner.Show()
}

func (l *roomViewLobby) hideSpinner() {
	l.spinner.Stop()
	l.spinner.Hide()
}

func (l *roomViewLobby) onJoinClicked() {
	l.disableFields()
	l.showSpinner()
	l.joinButton.SetSensitive(false)

	nickname, _ := l.nicknameEntry.GetText()
	password, _ := l.passwordEntry.GetText()

	l.onJoin = make(chan bool)
	l.onJoinError = make(chan error)
	l.cancel = make(chan bool)

	go l.sendJoinRoomRequest(nickname, password)
	go l.whenEnterRequestHasBeenResolved(nickname)
}

var (
	errJoinRequestFailed    = errors.New("the request to join the room has failed")
	errJoinNoConnection     = errors.New("join request failed because maybe no connection available")
	errJoinNicknameConflict = errors.New("join failed because the nickname is being used")
	errJoinOnlyMembers      = errors.New("join failed because only registered members are allowed")
	errJoinNotAuthorized    = errors.New("join failed because doesn't have authorization")
)

type mucRoomLobbyErr struct {
	room     jid.Bare
	nickname string
	errType  error
}

func (e *mucRoomLobbyErr) Error() string {
	return e.errType.Error()
}

func newMUCRoomLobbyErr(roomID jid.Bare, nickname string, errType error) error {
	return &mucRoomLobbyErr{
		room:     roomID,
		nickname: nickname,
		errType:  errType,
	}
}

func (l *roomViewLobby) joinRoom(nickname, password string) error {
	return l.account.session.JoinRoom(l.roomID, nickname, password)
}

func (l *roomViewLobby) sendJoinRoomRequest(nickname, password string) {
	err := l.joinRoom(nickname, password)
	if err != nil {
		l.log.WithField("nickname", nickname).WithError(err).Error("An error occurred while trying to join the room.")
		l.finishJoinRequest(errJoinNoConnection)
	}
}

func (l *roomViewLobby) whenEnterRequestHasBeenResolved(nickname string) {
	select {
	case <-l.onJoin:
		doInUIThread(func() {
			l.clearErrors()
			l.onSuccess()
		})
	case err := <-l.onJoinError:
		l.log.WithField("nickname", nickname).WithError(err).Error("An error occurred while trying to join the room")
		doInUIThread(func() {
			l.onJoinFailed(err)
		})
	case <-l.cancel:
	}
}

func (l *roomViewLobby) onJoinFailed(err error) {
	shouldEnableCreation := true

	userMessage := i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
	if err, ok := err.(*mucRoomLobbyErr); ok {
		userMessage = l.getUserErrorMessage(err)

		if err.errType == errJoinNicknameConflict {
			shouldEnableCreation = false
			if !l.nicknamesWithConflict.Has(err.nickname) {
				l.nicknamesWithConflict.Insert(err.nickname)
			}
		}
	}

	l.notifyOnError(userMessage)

	l.enableFields()
	l.hideSpinner()
	l.joinButton.SetSensitive(shouldEnableCreation)
}

func (l *roomViewLobby) getUserErrorMessage(err *mucRoomLobbyErr) string {
	switch err.errType {
	case errJoinNicknameConflict:
		return i18n.Local("Can't join the room using that nickname because it's already being used.")
	case errJoinOnlyMembers:
		return i18n.Local("Sorry, this room only allows registered members.")
	case errJoinNotAuthorized:
		return i18n.Local("Invalid password.")
	default:
		return i18n.Local("An error occurred while trying to join the room, please check your connection or make sure the room exists.")
	}
}

func (l *roomViewLobby) onJoinCancel() {
	if l.cancel != nil {
		l.cancel <- true
		l.cancel = nil
	}

	l.onCancel()
}

func (l *roomViewLobby) clearErrors() {
	l.errorNotif.Hide()
}

func (l *roomViewLobby) notifyOnError(err string) {
	if l.notification != nil {
		l.notificationArea.Remove(l.notification)
	}
	l.errorNotif.ShowMessage(err)
}

func (l *roomViewLobby) finishJoinRequest(err error) {
	if err != nil {
		l.onJoinError <- err
	} else {
		l.onJoin <- true
	}
}

func (l *roomViewLobby) occupantSelfJoinedEvent() {
	l.finishJoinRequest(nil)
}

func (l *roomViewLobby) joinFailed(roomID jid.Bare, nickname string) {
	l.finishJoinRequest(newMUCRoomLobbyErr(roomID, nickname, errJoinRequestFailed))
}

func (l *roomViewLobby) nicknameConflictEvent(roomID jid.Bare, nickname string) {
	l.finishJoinRequest(newMUCRoomLobbyErr(roomID, nickname, errJoinNicknameConflict))
}

func (l *roomViewLobby) registrationRequiredEvent(roomID jid.Bare, nickname string) {
	l.finishJoinRequest(newMUCRoomLobbyErr(roomID, nickname, errJoinOnlyMembers))
}

func (l *roomViewLobby) notAuthorizedEvent() {
	l.finishJoinRequest(newMUCRoomLobbyErr(nil, "", errJoinNotAuthorized))
}
