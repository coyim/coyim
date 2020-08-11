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

	account  *account
	jid      jid.Bare
	onCancel chan bool
	onJoin   chan bool

	connectionEventHandlers []func()

	events chan interface{}

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
		events:  make(chan interface{}),
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
	if rv.onCancel != nil {
		rv.onCancel <- true
		rv.onJoin <- false
	}

	doInUIThread(rv.clearErrors)

	rv.onCancel = make(chan bool, 1)
	rv.onJoin = make(chan bool, 1)
	nickName, _ := rv.nicknameEntry.GetText()

	doInUIThread(func() {
		rv.spinnerJoinView.Start()
		rv.spinnerJoinView.SetVisible(true)
	})

	go rv.account.session.JoinRoom(rv.jid, nickName)
	go func() {
		defer func() {
			doInUIThread(func() {
				rv.spinnerJoinView.Stop()
				rv.spinnerJoinView.SetVisible(false)
			})
		}()
		for {
			select {
			case jev, ok := <-rv.onJoin:
				if !ok {
					//TODO: Add the error message here
					return
				}
				if !jev {
					//TODO: Capture the error details here to show to the user
					rv.notifyOnError(i18n.Localf("You can't logged in to the room. Details: %s", jev))
					rv.account.log.WithField("Join Event: ", jev).Debug("Some user can't logged in to the room.")
				} else {
					rv.clearErrors()
					rv.togglePanelView()
				}
				return
			case _, _ = <-rv.onCancel:
				return
			}
		}
	}()
	go func() {
		//TODO: this event need to receive a data for the MUC Event received
		rv.onPresenceReceived(func() {
			rv.account.log.WithField("Join Channel: ", rv.onJoin).Info("Presence received...")
			rv.onJoin <- true
		})
	}()
}

func (rv *roomView) onPresenceReceived(f func()) {
	if rv.connectionEventHandlers == nil {
		rv.connectionEventHandlers = []func(){}
	}
	rv.connectionEventHandlers = append(rv.connectionEventHandlers, f)
}

func (u *gtkUI) viewForRoom(room *muc.Room) *roomView {
	if room.Opaque == nil {
		panic("developer error: trying to get an undefined view from room")
	}

	view, succeed := room.Opaque.(*roomView)
	if !succeed {
		panic("developer error: failed parsing room view into room.Opaque")
	}

	return view
}

func (u *gtkUI) mucShowRoom(a *account, ident jid.Bare) {
	room, err := a.joinRoom(u, ident)
	if err != nil {
		// TODO: Notify in a proper way this error
		log.Fatal(err.Error())
		return
	}
	view := u.viewForRoom(room)
	view.init()

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window":     func() {},
		"on_show_window":      view.validateInput,
		"on_nickname_changed": view.validateInput,
		"on_password_changed": view.validateInput,
		"on_password_checked": func() {
			view.togglePassword()
			view.validateInput()
		},
		"on_room_cancel_clicked": view.window.Destroy,
		"on_room_join_clicked": func() {
			view.onRoomJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.window)

	view.window.Show()
}
