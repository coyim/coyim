package gui

import (
	"sync"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
	log "github.com/sirupsen/logrus"
)

type roomView struct {
	builder *builder
	u       *gtkUI

	account          *account
	jid              jid.Bare
	onJoin           chan bool
	lastError        error
	lastErrorMessage string
	sync.RWMutex

	window           gtki.Window      `gtk-widget:"roomWindow"`
	boxJoinRoomView  gtki.Box         `gtk-widget:"boxJoinRoomView"`
	nicknameEntry    gtki.Entry       `gtk-widget:"nicknameEntry"`
	passwordCheck    gtki.CheckButton `gtk-widget:"passwordCheck"`
	passwordLabel    gtki.Label       `gtk-widget:"passwordLabel"`
	passwordEntry    gtki.Entry       `gtk-widget:"passwordEntry"`
	roomJoinButton   gtki.Button      `gtk-widget:"roomJoinButton"`
	spinnerJoinView  gtki.Spinner     `gtk-widget:"joinSpinner"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`
	boxRoomView      gtki.Box         `gtk-widget:"boxRoomView"`

	notification gtki.InfoBar
	errorNotif   *errorNotification
}

func (v *roomView) clearErrors() {
	v.errorNotif.Hide()
}

func (v *roomView) notifyOnError(err string) {
	if v.notification != nil {
		v.notificationArea.Remove(v.notification)
	}

	v.errorNotif.ShowMessage(err)
}

func (v *roomView) initUIBuilder() {
	v.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(v.builder.bindObjects(v))

	v.builder.ConnectSignals(map[string]interface{}{
		"on_show_window":         v.validateInput,
		"on_nickname_changed":    v.validateInput,
		"on_password_changed":    v.validateInput,
		"on_password_checked":    v.onPasswordChecked,
		"on_room_cancel_clicked": v.window.Destroy,
		"on_room_join_clicked":   v.joinRoom,
		"on_close_window":        v.onCloseWindow,
	})
}

func (v *roomView) onPasswordChecked() {
	v.setPasswordSensitiveBasedOnCheck()
	v.validateInput()
}

func (v *roomView) onCloseWindow() {
	_ = v.account.roomManager.LeaveRoom(v.jid)
}

func (v *roomView) initDefaults() {
	v.errorNotif = newErrorNotification(v.notificationArea)
	v.setPasswordSensitiveBasedOnCheck()
	v.window.SetTitle(i18n.Localf("Room: [%s]", v.jid))
}

func (v *roomView) setPasswordSensitiveBasedOnCheck() {
	a := v.passwordCheck.GetActive()
	v.passwordLabel.SetSensitive(a)
	v.passwordEntry.SetSensitive(a)
}

func (v *roomView) hasValidNickname() bool {
	nickname, _ := v.nicknameEntry.GetText()
	return len(nickname) > 0
}

func (v *roomView) hasValidPassword() bool {
	cv := v.passwordCheck.GetActive()
	if !cv {
		return true
	}
	password, _ := v.passwordEntry.GetText()
	return len(password) > 0
}

func (v *roomView) validateInput() {
	sensitiveValue := v.hasValidNickname() && v.hasValidPassword()
	v.roomJoinButton.SetSensitive(sensitiveValue)
}

func (v *roomView) togglePanelView() {
	doInUIThread(func() {
		value := v.boxJoinRoomView.IsVisible()
		v.boxJoinRoomView.SetVisible(!value)
		v.boxRoomView.SetVisible(value)
	})
}

func (v *roomView) startSpinner() {
	v.spinnerJoinView.Start()
	v.spinnerJoinView.SetVisible(true)
	v.roomJoinButton.SetSensitive(false)
}

func (v *roomView) stopSpinner() {
	v.spinnerJoinView.Stop()
	v.spinnerJoinView.SetVisible(false)
	v.roomJoinButton.SetSensitive(true)
}

func (v *roomView) joinRoom() {
	v.clearErrors()

	v.onJoin = make(chan bool)
	nickname, _ := v.nicknameEntry.GetText()

	v.startSpinner()

	go func() {
		err := v.account.session.JoinRoom(v.jid, nickname)
		if err != nil {
			doInUIThread(func() {
				v.stopSpinner()
				v.account.log.WithError(err).Error("Trying to join a room")
			})
		}
	}()

	go func() {
		defer func() {
			doInUIThread(func() {
				v.stopSpinner()
			})
		}()

		jev, ok := <-v.onJoin
		if !ok {
			//TODO: Add the error message here
			return
		}

		if !jev {
			doInUIThread(func() {
				v.notifyOnError(v.lastErrorMessage)
				v.account.log.WithFields(log.Fields{
					"room":     v.jid,
					"nickname": nickname,
					"message":  v.lastErrorMessage,
				}).Error("The user couldn't join a room")
			})
		} else {
			doInUIThread(func() {
				v.clearErrors()
				v.togglePanelView()
			})
		}
	}()
}

func (u *gtkUI) newRoom(a *account, ident jid.Bare) *muc.Room {
	room := muc.NewRoom(ident)

	view := &roomView{
		account: a,
		jid:     ident,
		u:       u,
	}

	view.initUIBuilder()
	view.initDefaults()

	room.Opaque = view

	return room

}

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare) {
	room, ok := a.roomManager.GetRoom(ident)
	if !ok {
		room = u.newRoom(a, ident)
		a.roomManager.AddRoom(room)
	}

	view := getViewFromRoom(room)

	if !ok {
		view.window.Show()
		return
	}

	view.window.Present()
}

func getViewFromRoom(r *muc.Room) *roomView {
	return r.Opaque.(*roomView)
}
