package gui

import (
	"log"
	"sync"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
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
	rv.togglePassword()

	rv.window.SetTitle(i18n.Localf("Room: [%s]", rv.jid))
}

// TODO[OB]-MUC: I don't think this is a great name for the method - since the method doesn't do what the method name says it does

func (rv *roomView) togglePassword() {
	doInUIThread(func() {
		value := rv.passwordCheck.GetActive()
		rv.passwordLabel.SetSensitive(value)
		rv.passwordEntry.SetSensitive(value)
	})
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

func (rv *roomView) onRoomJoinClicked() {
	doInUIThread(rv.clearErrors)

	// TODO[OB]-MUC: Hmm, this channel should likely not be cached. That seems to imply a race condition later on
	rv.onJoin = make(chan bool, 1)
	nickName, _ := rv.nicknameEntry.GetText()

	doInUIThread(func() {
		rv.spinnerJoinView.Start()
		rv.spinnerJoinView.SetVisible(true)
		rv.roomJoinButton.SetSensitive(false)
	})

	go rv.account.session.JoinRoom(rv.jid, nickName)
	go func() {
		defer func() {
			doInUIThread(func() {
				rv.spinnerJoinView.Stop()
				rv.spinnerJoinView.SetVisible(false)
				rv.roomJoinButton.SetSensitive(true)
			})
		}()

		jev, ok := <-rv.onJoin
		if !ok {
			//TODO: Add the error message here
			return
		}

		if !jev {
			// TODO[OB]-MUC: You should only send fixed strings to the i18n functions
			rv.notifyOnError(i18n.Local(rv.lastErrorMessage))
			// TODO[OB]-MUC: Why debug-level?
			// TODO[OB]-MUC: I don't think this log message makes much sense
			rv.account.log.WithField("Join Event: ", jev).Debug("Some user can't logged in to the room.")
		} else {
			rv.clearErrors()
			rv.togglePanelView()
		}
	}()
}

// TODO[OB]-MUC: This method name is confusing

func (rv *roomView) processOccupantJoinedEvent(err error) {
	// TODO[OB]-MUC: Why does this method take an error argument?

	if err != nil {
		// TODO[OB]-MUC: Why debug level?
		rv.account.log.WithError(err).Debug("Room join event received")
		rv.lastErrorMessage = err.Error()
		rv.onJoin <- false
		return
	}
	rv.onJoin <- true
}

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare) {
	_, err := u.addRoom(a, ident)
	if err != nil {
		// TODO: Notify in a proper way this error
		// TODO[OB]-MUC: Use log on 'a' object
		log.Fatal(err.Error())
		return
	}
	view, _, _ := a.roomViewFor(ident)
	view.init()

	view.builder.ConnectSignals(map[string]interface{}{
		// TODO[OB]-MUC: Useless handler
		"on_close_window":     func() {},
		"on_show_window":      view.validateInput,
		"on_nickname_changed": view.validateInput,
		"on_password_changed": view.validateInput,
		"on_password_checked": func() {
			view.togglePassword()
			view.validateInput()
		},
		"on_room_cancel_clicked": view.window.Destroy,
		// TODO[OB]-MUC: Can be inlined without func()
		"on_room_join_clicked": func() {
			view.onRoomJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.window)

	view.window.Show()
}
