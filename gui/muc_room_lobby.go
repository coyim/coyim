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
	ident jid.Bare
	ac    *account
	log   coylog.Logger

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

	onJoinChannel      chan bool
	onJoinErrorChannel chan error

	isReadyToJoinRoom bool

	nickNamesWithConflict *set.Set
	warnings              []*roomLobbyWarning

	// onSuccess will NOT be called from the UI thread
	onSuccess func()

	// onCancel will BE called from the UI thread
	onCancel func()
}

func (v *roomView) initRoomLobby() {
	v.lobby = v.newRoomViewLobby(v.account, v.identity, v.content, v.onJoined, v.onJoinCancel)
}

func (v *roomView) newRoomViewLobby(a *account, rid jid.Bare, parent gtki.Box, onSuccess, onCancel func()) *roomViewLobby {
	e := &roomViewLobby{
		ident:                 rid,
		ac:                    a,
		parent:                parent,
		onSuccess:             onSuccess,
		onCancel:              onCancel,
		nickNamesWithConflict: set.New(),
	}

	builder := newBuilder("MUCRoomLobby")
	panicOnDevError(builder.bindObjects(e))

	e.errorNotif = newErrorNotification(e.notificationArea)
	e.log = a.log.WithField("room", e.ident)

	builder.ConnectSignals(map[string]interface{}{
		"on_nickname_changed": e.onNickNameChange,
		"on_joined_clicked":   e.onJoin,
		"on_cancel_clicked":   e.onJoinCancel,
	})

	e.roomNameLabel.SetText(e.ident.String())
	e.content.SetHExpand(true)
	e.parent.Add(e.content)

	e.content.SetCenterWidget(e.mainBox)

	e.withRoomInfo(v.info)

	v.subscribe("lobby", occupantSelfJoined, e.occupantJoined)
	v.subscribe("lobby", roomInfoReceived, func(roomViewEventInfo) {
		e.withRoomInfo(v.info)
	})

	v.subscribe("lobby", nicknameConflict, func(ei roomViewEventInfo) {
		e.nicknameConflict(v.identity, ei.nickname)
	})

	v.subscribe("lobby", registrationRequired, func(ei roomViewEventInfo) {
		e.registrationRequired(v.identity, ei.nickname)
	})

	v.subscribe("lobby", previousToSwitchToMain, func(roomViewEventInfo) {
		v.unsubscribe("lobby", occupantSelfJoined)
		v.unsubscribe("lobby", roomInfoReceived)
	})

	return e
}

func (v *roomViewLobby) withRoomInfo(info *muc.RoomListing) {
	if info != nil {
		v.showRoomWarnings(info)
	}
	v.isReadyToJoinRoom = true
	v.enableJoinIfConditionsAreMet()
}

func (v *roomViewLobby) showRoomWarnings(roomInfo *muc.RoomListing) {
	v.clearWarnings()

	v.addWarning(i18n.Local("Please be aware that communication in chat rooms is not encrypted - anyone that can intercept communication between you and the server - and the server itself - will be able to see what you are saying in this chat room. Only join this room and communicate here if you trust the server to not be hostile."))

	switch roomInfo.Anonymity {
	case "semi":
		v.addWarning(i18n.Local("This room is partially anonymous. This means that only moderators can connect your nickname with your real username (your JID)."))
	case "no":
		v.addWarning(i18n.Local("This room is not anonymous. This means that any person in this room can connect your nickname with your real username (your JID)."))
	default:
		v.log.WithField("anonymity", roomInfo.Anonymity).Warn("Unknown anonymity setting for room")
	}

	if roomInfo.Logged {
		v.addWarning(i18n.Local("This room is publicly logged, meaning that everything you and the others in the room say or do can be made public on a website."))
	}
}

type roomLobbyWarning struct {
	text string

	bar     gtki.Box   `gtk-widget:"warning-infobar"`
	message gtki.Label `gtk-widget:"message"`
}

// addWarning should be called from the UI thread
func (v *roomViewLobby) addWarning(s string) {
	w := &roomLobbyWarning{text: s}
	v.warnings = append(v.warnings, w)

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

	v.warningsArea.PackStart(w.bar, false, false, 5)

	v.warningsArea.ShowAll()
}

func (v *roomViewLobby) clearWarnings() {
	for _, w := range v.warnings {
		v.warningsArea.Remove(w.bar)
	}
}

func (v *roomViewLobby) swtichToReturnOnCancel() {
	v.cancelButton.SetProperty("label", i18n.Local("Return"))
}

func (v *roomViewLobby) swtichToCancel() {
	v.cancelButton.SetProperty("label", i18n.Local("Cancel"))
}

func (v *roomViewLobby) show() {
	v.content.Show()
}

func (v *roomViewLobby) hide() {
	v.content.Hide()
}

func (v *roomViewLobby) close() {
	v.hide()
	v.parent.Remove(v.content)
}

func (v *roomViewLobby) onNickNameChange() {
	v.enableJoinIfConditionsAreMet()
}

func (v *roomViewLobby) enableJoinIfConditionsAreMet() {
	v.clearErrors()

	nickName, _ := v.nickNameEntry.GetText()
	conditionsAreValid := v.isReadyToJoinRoom && len(nickName) != 0

	if v.nickNamesWithConflict.Has(nickName) {
		conditionsAreValid = false
		v.notifyOnError(i18n.Local("That nickname is already being used."))
	}

	v.joinButton.SetSensitive(conditionsAreValid)
}

func (v *roomViewLobby) disableFields() {
	v.nickNameEntry.SetSensitive(false)
}

func (v *roomViewLobby) enableFields() {
	v.nickNameEntry.SetSensitive(true)
}

func (v *roomViewLobby) showSpinner() {
	v.spinner.Start()
	v.spinner.Show()
}

func (v *roomViewLobby) hideSpinner() {
	v.spinner.Stop()
	v.spinner.Hide()
}

func (v *roomViewLobby) onJoin() {
	v.disableFields()
	v.showSpinner()
	v.joinButton.SetSensitive(false)

	nickname, _ := v.nickNameEntry.GetText()

	v.onJoinChannel = make(chan bool)
	v.onJoinErrorChannel = make(chan error)

	go v.sendRoomEnterRequest(nickname)
	go v.whenEnterRequestHasBeenResolved(nickname)
}

var (
	errJoinRequestFailed    = errors.New("the request to join the room has failed")
	errJoinNoConnection     = errors.New("join request failed because maybe no connection available")
	errJoinNickNameConflict = errors.New("join failed because the nickname is being used")
	errJoinOnlyMembers      = errors.New("join failed because only registered members are allowed")
)

type mucRoomLobbyErr struct {
	room     jid.Bare
	nickname string
	errType  error
}

func (e *mucRoomLobbyErr) Error() string {
	return e.errType.Error()
}

func newMUCRoomLobbyErr(room jid.Bare, nickname string, errType error) error {
	return &mucRoomLobbyErr{
		room:     room,
		nickname: nickname,
		errType:  errType,
	}
}

func (v *roomViewLobby) sendRoomEnterRequest(nickName string) {
	err := v.ac.session.JoinRoom(v.ident, nickName)
	if err != nil {
		v.log.WithField("nickname", nickName).WithError(err).Error("An error occurred while trying to join the room.")
		v.finishJoinRequest(errJoinNoConnection)
	}
}

func (v *roomViewLobby) whenEnterRequestHasBeenResolved(nickname string) {
	select {
	case <-v.onJoinChannel:
		doInUIThread(v.clearErrors)
		if v.onSuccess != nil {
			v.onSuccess()
		}
	case err := <-v.onJoinErrorChannel:
		l := v.log.WithField("nickname", nickname)
		l.WithError(err).Error("An error occurred while trying to join the room")
		doInUIThread(func() {
			v.onJoinFailed(err)
		})
	}
}

func (v *roomViewLobby) onJoinFailed(err error) {
	shouldEnableCreation := true

	userMessage := i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
	if err, ok := err.(*mucRoomLobbyErr); ok {
		userMessage = v.getUserErrorMessage(err)

		if err.errType == errJoinNickNameConflict {
			shouldEnableCreation = false
			nickName := err.nickname
			if !v.nickNamesWithConflict.Has(nickName) {
				v.nickNamesWithConflict.Insert(nickName)
			}
		}
	}

	v.notifyOnError(userMessage)

	v.enableFields()
	v.hideSpinner()
	v.joinButton.SetSensitive(shouldEnableCreation)
}

func (v *roomViewLobby) getUserErrorMessage(err *mucRoomLobbyErr) string {
	switch err.errType {
	case errJoinNickNameConflict:
		return i18n.Local("Can't join the room using that nickname because it's already being used.")
	case errJoinOnlyMembers:
		return i18n.Local("Sorry, this room only allows registered members")
	default:
		return i18n.Local("An error occurred while trying to join the room, please check your connection or make sure the room exists.")
	}
}

func (v *roomViewLobby) onJoinCancel() {
	if v.onCancel != nil {
		v.onCancel()
	}
}

func (v *roomViewLobby) clearErrors() {
	v.errorNotif.Hide()
}

func (v *roomViewLobby) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}
	v.errorNotif.ShowMessage(err)
}

func (v *roomViewLobby) finishJoinRequest(err error) {
	if err != nil {
		v.onJoinErrorChannel <- err
	} else {
		v.onJoinChannel <- true
	}
}

func (v *roomViewLobby) occupantJoined(roomViewEventInfo) {
	v.finishJoinRequest(nil)
}

func (v *roomViewLobby) joinFailed(room jid.Bare, nickname string) {
	v.finishJoinRequest(newMUCRoomLobbyErr(room, nickname, errJoinRequestFailed))
}

func (v *roomViewLobby) nicknameConflict(room jid.Bare, nickname string) {
	v.finishJoinRequest(newMUCRoomLobbyErr(room, nickname, errJoinNickNameConflict))
}

func (v *roomViewLobby) registrationRequired(room jid.Bare, nickname string) {
	v.finishJoinRequest(newMUCRoomLobbyErr(room, nickname, errJoinOnlyMembers))
}
