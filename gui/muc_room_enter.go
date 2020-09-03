package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomViewEnter struct {
	ident jid.Bare
	ac    *account
	log   coylog.Logger

	content          gtki.Box         `gtk-widget:"boxJoinRoomView"`
	roomNameLabel    gtki.Label       `gtk-widget:"roomNameValue"`
	nickNameEntry    gtki.Entry       `gtk-widget:"nickNameEntry"`
	passwordEntry    gtki.Entry       `gtk-widget:"passwordEntry"`
	passwordCheck    gtki.CheckButton `gtk-widget:"passwordCheck"`
	enterButton      gtki.Button      `gtk-widget:"enterButton"`
	spinner          gtki.Spinner     `gtk-widget:"spinner"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`
	parent           gtki.Box
	errorNotif       *errorNotification
	notification     gtki.InfoBar

	lastError        error
	lastErrorMessage string
	onEnterChannel   chan bool

	onSuccess func()
	onCancel  func()
}

func newRoomEnterView(a *account, rid jid.Bare, parent gtki.Box, onSuccess, onCancel func()) *roomViewEnter {
	e := &roomViewEnter{
		ident:     rid,
		ac:        a,
		parent:    parent,
		onSuccess: onSuccess,
		onCancel:  onCancel,
	}

	builder := newBuilder("MUCRoomJoin")
	panicOnDevError(builder.bindObjects(e))

	e.errorNotif = newErrorNotification(e.notificationArea)
	e.log = a.log.WithField("room", e.ident)

	builder.ConnectSignals(map[string]interface{}{
		"on_nickname_changed": e.onNickNameChange,
		"on_password_changed": e.onPasswordChange,
		"on_password_checked": e.onUsePasswordChanged,
		"on_enter_clicked":    e.onEnter,
		"on_cancel_clicked":   e.onEnterCancel,
	})

	e.roomNameLabel.SetText(e.ident.String())
	e.content.SetHExpand(true)
	e.parent.Add(e.content)

	return e
}

func (v *roomViewEnter) show() {
	v.content.Show()
}

func (v *roomViewEnter) hide() {
	v.content.Hide()
}

func (v *roomViewEnter) close() {
	v.hide()
	v.parent.Remove(v.content)
}

func (v *roomViewEnter) onNickNameChange() {
	v.enableJoinIfConditionsAreMet()
}

// TODO: Should we active the password checkbox
// if the user start to write a password?
func (v *roomViewEnter) onPasswordChange() {
	v.enableJoinIfConditionsAreMet()
}

func (v *roomViewEnter) onUsePasswordChanged() {
	isChecked := v.passwordCheck.GetActive()
	v.passwordEntry.SetSensitive(isChecked)
	if !isChecked {
		v.passwordEntry.SetText("")
	}
}

func (v *roomViewEnter) enableJoinIfConditionsAreMet() {
	nickName, _ := v.nickNameEntry.GetText()
	password, _ := v.passwordEntry.GetText()
	passwordIsActive := v.passwordCheck.GetActive()

	hasAllValues := len(nickName) != 0 && (!passwordIsActive || len(password) != 0)
	v.enterButton.SetSensitive(hasAllValues)
}

func (v *roomViewEnter) disableFields() {
	v.nickNameEntry.SetSensitive(false)
	v.passwordEntry.SetSensitive(false)
	v.passwordCheck.SetSensitive(false)
}

func (v *roomViewEnter) enableFields() {
	v.nickNameEntry.SetSensitive(true)
	v.passwordEntry.SetSensitive(true)
	v.passwordCheck.SetSensitive(true)
}

func (v *roomViewEnter) showSpinner() {
	v.spinner.Start()
	v.spinner.Show()
}

func (v *roomViewEnter) hideSpinner() {
	v.spinner.Stop()
	v.spinner.Hide()
}

func (v *roomViewEnter) onEnter() {
	v.disableFields()
	v.showSpinner()
	v.enterButton.SetSensitive(false)

	nickName, _ := v.nickNameEntry.GetText()

	v.onEnterChannel = make(chan bool)
	go v.sendRoomEnterRequest(nickName)
	go v.whenEnterRequestHasBeenResolved(nickName)
}

func (v *roomViewEnter) sendRoomEnterRequest(nickName string) {
	err := v.ac.session.JoinRoom(v.ident, nickName)
	if err != nil {
		v.log.WithField("nickname", nickName).WithError(err).Error("An error occurred while trying to join the room.")
		v.onEnterChannel <- false
		doInUIThread(v.onEnterFails)
	}
}

func (v *roomViewEnter) whenEnterRequestHasBeenResolved(nickName string) {
	hasJoined, ok := <-v.onEnterChannel
	if !ok {
		doInUIThread(func() {
			v.notifyOnError(i18n.Local("An error happened while trying to join the room, please check your connection or try again."))
		})
		return
	}

	if !hasJoined {
		if len(v.lastErrorMessage) == 0 {
			v.lastErrorMessage = i18n.Local("An error happened while trying to join the room, please check your connection or make sure the room exists.")
		}

		v.log.WithFields(log.Fields{
			"nickname": nickName,
			"message":  v.lastErrorMessage,
		}).Error("An error happened while trying to join the room")

		doInUIThread(func() {
			v.notifyOnError(v.lastErrorMessage)
		})

		return
	}

	doInUIThread(v.clearErrors)

	if v.onSuccess != nil {
		v.onSuccess()
	}
}

func (v *roomViewEnter) onEnterFails() {
	v.enableFields()
	v.hideSpinner()
	v.enterButton.SetSensitive(true)
}

func (v *roomViewEnter) onEnterCancel() {
	if v.onCancel != nil {
		v.onCancel()
	}
}

func (v *roomViewEnter) clearErrors() {
	v.errorNotif.Hide()
}

func (v *roomViewEnter) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}
	v.errorNotif.ShowMessage(err)
}

func (v *roomViewEnter) onRoomOccupantJoinedReceived() {
	v.onEnterChannel <- true
}

func (v *roomViewEnter) onJoinErrorRecevied(from jid.Full) {
	v.onEnterChannel <- false
}

func (v *roomViewEnter) onNicknameConflictReceived(from jid.Full) {
	v.lastErrorMessage = i18n.Localf("Can't join the room using \"%s\" because the nickname is already being used.", from.Resource())
	v.onEnterChannel <- false
}

func (v *roomViewEnter) onRegistrationRequiredReceived(from jid.Full) {
	v.lastErrorMessage = i18n.Local("Sorry, this room only allows registered members")
	v.onEnterChannel <- false
}
