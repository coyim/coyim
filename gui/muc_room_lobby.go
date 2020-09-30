package gui

import (
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
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
	nickNameEntry    gtki.Entry   `gtk-widget:"nickNameEntry"`
	joinButton       gtki.Button  `gtk-widget:"joinButton"`
	cancelButton     gtki.Button  `gtk-widget:"cancelButton"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"boxNotificationArea"`
	warningsArea     gtki.Box     `gtk-widget:"boxWarningsArea"`
	parent           gtki.Box
	errorNotif       *errorNotification
	notification     gtki.InfoBar

	onJoin      chan bool
	onJoinError chan error

	isReadyToJoinRoom bool

	nicknamesWithConflict *set.Set
	warnings              []*roomLobbyWarning

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
		log:                   v.log,
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
	l.initDefaults(v)
	l.initSubscribers(v)

	return l
}

func (l *roomViewLobby) initBuilder() {
	builder := newBuilder("MUCRoomLobby")
	panicOnDevError(builder.bindObjects(l))

	builder.ConnectSignals(map[string]interface{}{
		"on_nickname_changed": l.onNickNameChange,
		"on_joined_clicked":   l.onJoinClicked,
		"on_cancel_clicked":   l.onJoinCancel,
	})
}

func (l *roomViewLobby) initDefaults(v *roomView) {
	l.errorNotif = newErrorNotification(l.notificationArea)

	l.roomNameLabel.SetText(l.roomID.String())
	l.content.SetHExpand(true)
	l.parent.Add(l.content)
	l.content.SetCenterWidget(l.mainBox)

	l.setRoomInfo(v.info)
}

func (l *roomViewLobby) initSubscribers(v *roomView) {
	v.subscribeAll("lobby", roomViewEventObservers{
		"occupantSelfJoinedEvent": l.occupantJoinedEvent,
		"roomInfoReceivedEvent": func(roomViewEventInfo) {
			doInUIThread(func() {
				l.setRoomInfo(v.info)
			})
		},
		"nicknameConflictEvent": func(ei roomViewEventInfo) {
			doInUIThread(func() {
				l.nicknameConflictEvent(v.roomID(), ei["nickname"])
			})
		},
		"registrationRequiredEvent": func(ei roomViewEventInfo) {
			doInUIThread(func() {
				l.registrationRequiredEvent(v.roomID(), ei["nickname"])
			})
		},
		"beforeSwitchingToMainViewEvent": func(roomViewEventInfo) {
			v.unsubscribe("lobby", "occupantSelfJoinedEvent")
			v.unsubscribe("lobby", "roomInfoReceivedEvent")
		},
	})
}

func (l *roomViewLobby) setRoomInfo(info *muc.RoomListing) {
	l.clearWarnings()

	if info != nil {
		l.showRoomWarnings(info)
	}

	l.isReadyToJoinRoom = true
	l.enableJoinIfConditionsAreMet()
}

func (l *roomViewLobby) showRoomWarnings(roomInfo *muc.RoomListing) {
	l.addWarning(i18n.Local("Please be aware that communication in chat rooms is not encrypted - anyone that can intercept communication between you and the server - and the server itself - will be able to see what you are saying in this chat room. Only join this room and communicate here if you trust the server to not be hostile."))

	switch roomInfo.Anonymity {
	case "semi":
		l.addWarning(i18n.Local("This room is partially anonymous. This means that only moderators can connect your nickname with your real username (your JID)."))
	case "no":
		l.addWarning(i18n.Local("This room is not anonymous. This means that any person in this room can connect your nickname with your real username (your JID)."))
	default:
		l.log.WithField("anonymity", roomInfo.Anonymity).Warn("Unknown anonymity setting for room")
	}

	if roomInfo.Logged {
		l.addWarning(i18n.Local("This room is publicly logged, meaning that everything you and the others in the room say or do can be made public on a website."))
	}
}

type roomLobbyWarning struct {
	text string

	bar     gtki.Box   `gtk-widget:"warning-infobar"`
	message gtki.Label `gtk-widget:"message"`
}

// addWarning should be called from the UI thread
func (l *roomViewLobby) addWarning(s string) {
	w := &roomLobbyWarning{text: s}
	l.warnings = append(l.warnings, w)

	builder := newBuilder("MUCRoomWarning")
	panicOnDevError(builder.bindObjects(w))

	w.message.SetText(w.text)

	prov := providerWithStyle("box", style{
		"color":            "#744210",
		"background-color": "#fefcbf",
		"border":           "1px solid #d69e2e",
		"border-radius":    "4px",
		"padding":          "10px",
	})

	updateWithStyle(w.bar, prov)

	l.warningsArea.PackStart(w.bar, false, false, 5)

	l.warningsArea.ShowAll()
}

func (l *roomViewLobby) clearWarnings() {
	// TODO: Why can't we just remove
	// all the entitites inside the warningsArea
	// and then remove the need to have the "warnings" field at all
	for _, w := range l.warnings {
		l.warningsArea.Remove(w.bar)
	}
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

func (l *roomViewLobby) onNickNameChange() {
	l.enableJoinIfConditionsAreMet()
}

func (l *roomViewLobby) enableJoinIfConditionsAreMet() {
	l.clearErrors()

	nickname, _ := l.nickNameEntry.GetText()
	conditionsAreValid := l.isReadyToJoinRoom && nickname != ""

	if l.nicknamesWithConflict.Has(nickname) {
		conditionsAreValid = false
		l.notifyOnError(i18n.Local("That nickname is already being used."))
	}

	l.joinButton.SetSensitive(conditionsAreValid)
}

func (l *roomViewLobby) disableFields() {
	l.nickNameEntry.SetSensitive(false)
}

func (l *roomViewLobby) enableFields() {
	l.nickNameEntry.SetSensitive(true)
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

	nickname, _ := l.nickNameEntry.GetText()

	l.onJoin = make(chan bool)
	l.onJoinError = make(chan error)

	go l.sendRoomEnterRequest(nickname)
	go l.whenEnterRequestHasBeenResolved(nickname)
}

var (
	errJoinRequestFailed         = errors.New("the request to join the room has failed")
	errJoinNoConnection          = errors.New("join request failed because maybe no connection available")
	errJoinnicknameConflictEvent = errors.New("join failed because the nickname is being used")
	errJoinOnlyMembers           = errors.New("join failed because only registered members are allowed")
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

func (l *roomViewLobby) sendRoomEnterRequest(nickname string) {
	err := l.account.session.JoinRoom(l.roomID, nickname)
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
	}
}

func (l *roomViewLobby) onJoinFailed(err error) {
	shouldEnableCreation := true

	userMessage := i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
	if err, ok := err.(*mucRoomLobbyErr); ok {
		userMessage = l.getUserErrorMessage(err)

		if err.errType == errJoinnicknameConflictEvent {
			shouldEnableCreation = false
			// TODO: This naming is different from most other places
			nickName := err.nickname
			if !l.nicknamesWithConflict.Has(nickName) {
				l.nicknamesWithConflict.Insert(nickName)
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
	case errJoinnicknameConflictEvent:
		return i18n.Local("Can't join the room using that nickname because it's already being used.")
	case errJoinOnlyMembers:
		return i18n.Local("Sorry, this room only allows registered members.")
	default:
		return i18n.Local("An error occurred while trying to join the room, please check your connection or make sure the room exists.")
	}
}

func (l *roomViewLobby) onJoinCancel() {
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

func (l *roomViewLobby) occupantJoinedEvent(roomViewEventInfo) {
	l.finishJoinRequest(nil)
}

func (l *roomViewLobby) joinFailed(roomID jid.Bare, nickname string) {
	l.finishJoinRequest(newMUCRoomLobbyErr(roomID, nickname, errJoinRequestFailed))
}

func (l *roomViewLobby) nicknameConflictEvent(roomID jid.Bare, nickname string) {
	l.finishJoinRequest(newMUCRoomLobbyErr(roomID, nickname, errJoinnicknameConflictEvent))
}

func (l *roomViewLobby) registrationRequiredEvent(roomID jid.Bare, nickname string) {
	l.finishJoinRequest(newMUCRoomLobbyErr(roomID, nickname, errJoinOnlyMembers))
}
