package gui

import (
	"log"
	"sync"

	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/access"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomView struct {
	room    *muc.Room
	info    *muc.RoomListing
	account *account
	jid     jid.Bare

	session             access.Session
	sessionObserver     chan interface{}
	sessionObserverLock sync.RWMutex

	log coylog.Logger

	events chan interface{}

	sync.RWMutex

	builder *builder
	u       *gtkUI

	window           gtki.Window      `gtk-widget:"room-window"`
	boxJoinRoomView  gtki.Box         `gtk-widget:"boxJoinRoomView"`
	textNickname     gtki.Entry       `gtk-widget:"textNickname"`
	chkPassword      gtki.CheckButton `gtk-widget:"checkPassword"`
	labelPassword    gtki.Label       `gtk-widget:"labelPassword"`
	textPassword     gtki.Entry       `gtk-widget:"textPassword"`
	btnAcceptJoin    gtki.Button      `gtk-widget:"btnAcceptJoin"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification

	boxRoomView gtki.Box `gtk-widget:"boxRoomView"`
}

func newRoom(a *account, ident jid.Bare) *roomView {
	r := &roomView{
		room:    muc.NewRoom(ident),
		account: a,
		jid:     ident,
	}
	return r
}

func (r *roomView) init() {
	r.builder = newBuilder("MUCRoomWindow")

	panicOnDevError(r.builder.bindObjects(r))

	r.errorNotif = newErrorNotification(r.notificationArea)
	r.togglePassword()

	r.window.SetTitle(i18n.Localf("Room: [%s]", r.jid.String()))
}

func (r *roomView) togglePassword() {
	doInUIThread(func() {
		value := r.chkPassword.GetActive()
		r.labelPassword.SetSensitive(value)
		r.textPassword.SetSensitive(value)
	})
}

func (r *roomView) hasValidNickname() bool {
	nickName, _ := r.textNickname.GetText()
	return len(nickName) > 0
}

func (r *roomView) hasValidPassword() bool {
	value := r.chkPassword.GetActive()
	if !value {
		return true
	}
	password, _ := r.textPassword.GetText()
	return len(password) > 0
}

func (r *roomView) validateInput() {
	sensitiveValue := r.hasValidNickname() && r.hasValidPassword()
	r.btnAcceptJoin.SetSensitive(sensitiveValue)
}

func (r *roomView) togglePanelView() {
	doInUIThread(func() {
		value := r.boxJoinRoomView.IsVisible()
		r.boxJoinRoomView.SetVisible(!value)
		r.boxRoomView.SetVisible(value)
	})
}

func (r *roomView) onBtnJoinClicked() {
	nickName, _ := r.textNickname.GetText()
	go r.account.session.JoinRoom(r.jid, nickName)
	r.togglePanelView()
}

func (r *roomView) id() string {
	return r.room.Identity.String()
}

func (u *gtkUI) mucShowRoom(a *account, rjid jid.Bare) {
	view, err := a.joinRoom(u, rjid)
	if err != nil {
		// TODO: Notify in a proper way this error
		log.Fatal(err.Error())
		return
	}

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
		"on_cancel_join_clicked": view.window.Destroy,
		"on_accept_join_clicked": func() {
			view.onBtnJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.window)

	view.window.Show()
}
