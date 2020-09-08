package gui

import (
	"errors"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
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
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"boxNotificationArea"`
	warningsArea     gtki.Box     `gtk-widget:"boxWarningsArea"`
	parent           gtki.Box
	errorNotif       *errorNotification
	notification     gtki.InfoBar

	onJoinChannel      chan bool
	onJoinErrorChannel chan error

	onSuccess func()
	onCancel  func()
}

func newRoomViewLobby(a *account, rid jid.Bare, parent gtki.Box, onSuccess, onCancel func(), roomInfo *muc.RoomListing) *roomViewLobby {
	e := &roomViewLobby{
		ident:     rid,
		ac:        a,
		parent:    parent,
		onSuccess: onSuccess,
		onCancel:  onCancel,
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

	e.addWarning(i18n.Local("Please be aware that communication in chat rooms is not encrypted - anyone that can intercept communication between you and the server - and the server itself - will be able to see what you are saying in this chat room."))

	switch roomInfo.Anonymity {
	case "semi":
		e.addWarning(i18n.Local("This room is partially anonymous. This means that only moderators can connect your nickname with your real username (your JID)."))
	case "no":
		e.addWarning(i18n.Local("This room is not anomyous. This means that any person in this room can connect your nickname with your real username (your JID)."))
	default:
		e.log.WithField("anonymity", roomInfo.Anonymity).Warn("Unknown anonymity setting for room")
	}

	return e
}

type roomLobbyWarning struct {
	text string

	bar     gtki.Box   `gtk-widget:"warning-infobar"`
	message gtki.Label `gtk-widget:"message"`
}

// addWarning should be called from the UI thread
func (v *roomViewLobby) addWarning(s string) {
	w := &roomLobbyWarning{text: s}

	builder := newBuilder("MUCRoomWarning")
	panicOnDevError(builder.bindObjects(w))

	w.message.SetText(w.text)

	prov := providerWithCSS("box { background-color: #89AF8F; color: #000000; border: 1px solid #000000; border-radius: 5px; }")
	updateWithStyle(w.bar, prov)

	v.warningsArea.PackStart(w.bar, false, false, 5)

	v.warningsArea.ShowAll()
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
	nickName, _ := v.nickNameEntry.GetText()
	v.joinButton.SetSensitive(len(nickName) != 0)
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

	nickName, _ := v.nickNameEntry.GetText()

	v.onJoinChannel = make(chan bool)
	v.onJoinErrorChannel = make(chan error)

	go v.sendRoomEnterRequest(nickName)
	go v.whenEnterRequestHasBeenResolved(nickName)
}

var (
	errJoinRequestFailed    = errors.New("the request to join the room has failed")
	errJoinNoConnection     = errors.New("join request failed because maybe no connection available")
	errJoinNickNameConflict = errors.New("join failed because the nickname is being used")
	errJoinOnlyMembers      = errors.New("join failed because only registered members are allowed")
)

type mucRoomLobbyErr struct {
	from    jid.Full
	errType error
}

func (e *mucRoomLobbyErr) Error() string {
	return e.errType.Error()
}

func newMUCRoomLobbyErr(from jid.Full, errType error) error {
	return &mucRoomLobbyErr{
		from:    from,
		errType: errType,
	}
}

func (v *roomViewLobby) sendRoomEnterRequest(nickName string) {
	err := v.ac.session.JoinRoom(v.ident, nickName)
	if err != nil {
		v.log.WithField("nickname", nickName).WithError(err).Error("An error occurred while trying to join the room.")
		v.finishJoinRequest(errJoinNoConnection)
	}
}

func (v *roomViewLobby) whenEnterRequestHasBeenResolved(nickName string) {
	select {
	case <-v.onJoinChannel:
		doInUIThread(v.clearErrors)
		if v.onSuccess != nil {
			v.onSuccess()
		}
	case err := <-v.onJoinErrorChannel:
		l := v.log.WithField("nickname", nickName)

		l.WithError(err).Error("An error occurred while trying to join the room")
		doInUIThread(func() {
			v.onJoinFailed(err)
		})
	}
}

func (v *roomViewLobby) onJoinFailed(err error) {
	userMessage := v.getUserErrorMessage(err)
	v.notifyOnError(userMessage)

	v.enableFields()
	v.hideSpinner()
	v.joinButton.SetSensitive(true)
}

func (v *roomViewLobby) getUserErrorMessage(err error) string {
	if err, ok := err.(*mucRoomLobbyErr); ok {
		switch err.errType {
		case errJoinRequestFailed:
			return i18n.Local("An error occurred while trying to join the room, please check your connection or make sure the room exists.")
		case errJoinNickNameConflict:
			return i18n.Localf("Can't join the room using the nickname \"%s\" because it's already being used.", err.from.Resource())
		case errJoinOnlyMembers:
			return i18n.Local("Sorry, this room only allows registered members")
		}
	}
	return i18n.Local("An unknown error occurred while trying to join the room, please check your connection or try again.")
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

func (v *roomViewLobby) onRoomOccupantJoinedReceived() {
	v.finishJoinRequest(nil)
}

func (v *roomViewLobby) onJoinErrorRecevied(from jid.Full) {
	v.finishJoinRequest(newMUCRoomLobbyErr(from, errJoinRequestFailed))
}

func (v *roomViewLobby) onNicknameConflictReceived(from jid.Full) {
	v.finishJoinRequest(newMUCRoomLobbyErr(from, errJoinNickNameConflict))
}

func (v *roomViewLobby) onRegistrationRequiredReceived(from jid.Full) {
	v.finishJoinRequest(newMUCRoomLobbyErr(from, errJoinOnlyMembers))
}
