package gui

import (
	"errors"
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

	window           gtki.Window      `gtk-widget:"room-window"`
	boxJoinRoomView  gtki.Box         `gtk-widget:"boxJoinRoomView"`
	nicknameEntry    gtki.Entry       `gtk-widget:"nicknameEntry"`
	passwordCheck    gtki.CheckButton `gtk-widget:"passwordCheck"`
	passwordLabel    gtki.Label       `gtk-widget:"passwordLabel"`
	passwordEntry    gtki.Entry       `gtk-widget:"passwordEntry"`
	roomJoinButton   gtki.Button      `gtk-widget:"roomJoinButton"`
	spinnerJoinView  gtki.Spinner     `gtk-widget:"joinSpinner"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification

	boxRoomView gtki.Box `gtk-widget:"boxRoomView"`
}

func (rv *roomView) clearErrors() {
	rv.errorNotif.Hide()
}

func (rv *roomView) notifyOnError(err string) {
	if rv.notification != nil {
		rv.notificationArea.Remove(rv.notification)
	}

	rv.errorNotif.ShowMessage(err)
}

func (u *gtkUI) newRoom(a *account, ident jid.Bare) *muc.Room {
	r := muc.NewRoom(ident)
	r.Opaque = &roomView{
		account: a,
		jid:     ident,
		u:       u,
	}

	return r
}

func (rv *roomView) init() {
	rv.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(rv.builder.bindObjects(rv))

	rv.errorNotif = newErrorNotification(rv.notificationArea)
	rv.setPasswordSensitiveBasedOnCheck()

	rv.window.SetTitle(i18n.Localf("Room: [%s]", rv.jid))
}

func (rv *roomView) setPasswordSensitiveBasedOnCheck() {
	v := rv.passwordCheck.GetActive()

	rv.passwordLabel.SetSensitive(v)
	rv.passwordEntry.SetSensitive(v)
}

func (rv *roomView) hasValidNickname() bool {
	nickName, _ := rv.nicknameEntry.GetText()
	return len(nickName) > 0
}

func (rv *roomView) hasValidPassword() bool {
	cv := rv.passwordCheck.GetActive()
	if !cv {
		return true
	}
	password, _ := rv.passwordEntry.GetText()
	return len(password) > 0
}

func (rv *roomView) validateInput() {
	sensitiveValue := rv.hasValidNickname() && rv.hasValidPassword()
	rv.roomJoinButton.SetSensitive(sensitiveValue)
}

func (rv *roomView) togglePanelView() {
	doInUIThread(func() {
		value := rv.boxJoinRoomView.IsVisible()
		rv.boxJoinRoomView.SetVisible(!value)
		rv.boxRoomView.SetVisible(value)
	})
}

func (rv *roomView) startSpinner() {
	rv.spinnerJoinView.Start()
	rv.spinnerJoinView.SetVisible(true)
	rv.roomJoinButton.SetSensitive(false)
}

func (rv *roomView) stopSpinner() {
	rv.spinnerJoinView.Stop()
	rv.spinnerJoinView.SetVisible(false)
	rv.roomJoinButton.SetSensitive(true)
}

func (rv *roomView) onRoomJoinClicked() {
	rv.clearErrors()

	rv.onJoin = make(chan bool)
	nickName, _ := rv.nicknameEntry.GetText()

	rv.startSpinner()

	go func() {
		err := rv.account.session.JoinRoom(rv.jid, nickName)
		if err != nil {
			doInUIThread(func() {
				rv.stopSpinner()
				rv.account.log.WithError(err).Error("Trying to join a room")
			})
		}
	}()
	go func() {
		defer func() {
			doInUIThread(func() {
				rv.stopSpinner()
			})
		}()

		jev, ok := <-rv.onJoin
		if !ok {
			//TODO: Add the error message here
			return
		}

		if !jev {
			doInUIThread(func() {
				rv.notifyOnError(rv.lastErrorMessage)
				rv.account.log.WithFields(log.Fields{
					"Room":     rv.jid,
					"Nickname": nickName,
					"Message":  rv.lastErrorMessage,
				}).Error("The user couldn't join a room")
			})
		} else {
			doInUIThread(func() {
				rv.clearErrors()
				rv.togglePanelView()
			})
		}
	}()
}

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare) {
	room, err := u.addRoom(a, ident)
	if err != nil {
		// TODO: Notify in a proper way this error
		a.log.Fatal(err.Error())
		return
	}

	view, ok := room.Opaque.(*roomView)
	if !ok {
		// TODO: Notify in a proper way this error
		a.log.Fatal(errors.New("Can't create the room view"))
	}
	view.init()

	view.builder.ConnectSignals(map[string]interface{}{
		"on_show_window":      view.validateInput,
		"on_nickname_changed": view.validateInput,
		"on_password_changed": view.validateInput,
		"on_password_checked": func() {
			view.setPasswordSensitiveBasedOnCheck()
			view.validateInput()
		},
		"on_room_cancel_clicked": view.window.Destroy,
		"on_room_join_clicked":   view.onRoomJoinClicked,
	})

	u.connectShortcutsChildWindow(view.window)

	view.window.Show()
}
