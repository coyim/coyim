package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomViewJoin struct {
	ident jid.Bare
	ac    *account
	log   coylog.Logger

	content          gtki.Box     `gtk-widget:"boxJoinRoomView"`
	roomNameLabel    gtki.Label   `gtk-widget:"roomNameValue"`
	nickNameEntry    gtki.Entry   `gtk-widget:"nickNameEntry"`
	joinButton       gtki.Button  `gtk-widget:"joinButton"`
	spinner          gtki.Spinner `gtk-widget:"spinner"`
	notificationArea gtki.Box     `gtk-widget:"boxNotificationArea"`
	parent           gtki.Box
	errorNotif       *errorNotification
	notification     gtki.InfoBar

	lastError        error
	lastErrorMessage string
	onJoinChannel    chan bool

	onSuccess func()
	onCancel  func()
}

func newRoomEnterView(a *account, rid jid.Bare, parent gtki.Box, onSuccess, onCancel func()) *roomViewJoin {
	e := &roomViewJoin{
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
		"on_joined_clicked":   e.onJoin,
		"on_cancel_clicked":   e.onJoinCancel,
	})

	e.roomNameLabel.SetText(e.ident.String())
	e.content.SetHExpand(true)
	e.parent.Add(e.content)

	return e
}

func (v *roomViewJoin) show() {
	v.content.Show()
}

func (v *roomViewJoin) hide() {
	v.content.Hide()
}

func (v *roomViewJoin) close() {
	v.hide()
	v.parent.Remove(v.content)
}

func (v *roomViewJoin) onNickNameChange() {
	v.enableJoinIfConditionsAreMet()
}

func (v *roomViewJoin) enableJoinIfConditionsAreMet() {
	nickName, _ := v.nickNameEntry.GetText()
	v.joinButton.SetSensitive(len(nickName) != 0)
}

func (v *roomViewJoin) disableFields() {
	v.nickNameEntry.SetSensitive(false)
}

func (v *roomViewJoin) enableFields() {
	v.nickNameEntry.SetSensitive(true)
}

func (v *roomViewJoin) showSpinner() {
	v.spinner.Start()
	v.spinner.Show()
}

func (v *roomViewJoin) hideSpinner() {
	v.spinner.Stop()
	v.spinner.Hide()
}

func (v *roomViewJoin) onJoin() {
	v.disableFields()
	v.showSpinner()
	v.joinButton.SetSensitive(false)

	nickName, _ := v.nickNameEntry.GetText()

	v.onJoinChannel = make(chan bool)
	go v.sendRoomEnterRequest(nickName)
	go v.whenEnterRequestHasBeenResolved(nickName)
}

func (v *roomViewJoin) sendRoomEnterRequest(nickName string) {
	err := v.ac.session.JoinRoom(v.ident, nickName)
	if err != nil {
		v.log.WithField("nickname", nickName).WithError(err).Error("An error occurred while trying to join the room.")
		v.onJoinChannel <- false
		doInUIThread(v.onJoinFails)
	}
}

func (v *roomViewJoin) whenEnterRequestHasBeenResolved(nickName string) {
	hasJoined, ok := <-v.onJoinChannel
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

func (v *roomViewJoin) onJoinFails() {
	v.enableFields()
	v.hideSpinner()
	v.joinButton.SetSensitive(true)
}

func (v *roomViewJoin) onJoinCancel() {
	if v.onCancel != nil {
		v.onCancel()
	}
}

func (v *roomViewJoin) clearErrors() {
	v.errorNotif.Hide()
}

func (v *roomViewJoin) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}
	v.errorNotif.ShowMessage(err)
}

func (v *roomViewJoin) onRoomOccupantJoinedReceived() {
	v.onJoinChannel <- true
}

func (v *roomViewJoin) onJoinErrorRecevied(from jid.Full) {
	v.onJoinChannel <- false
}

func (v *roomViewJoin) onNicknameConflictReceived(from jid.Full) {
	v.lastErrorMessage = i18n.Localf("Can't join the room using \"%s\" because the nickname is already being used.", from.Resource())
	v.onJoinChannel <- false
}

func (v *roomViewJoin) onRegistrationRequiredReceived(from jid.Full) {
	v.lastErrorMessage = i18n.Local("Sorry, this room only allows registered members")
	v.onJoinChannel <- false
}
