package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomView struct {
	builder *builder

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

	// TODO: this is temporary.
	// the bellow fields should be an interface with all the necessary room information
	roomJid  jid.Bare
	roomInfo *muc.RoomListing
	account  *account
}

func (rv *mucRoomView) init() {
	rv.builder = newBuilder("MUCRoomWindow")
	panicOnDevError(rv.builder.bindObjects(rv))
	rv.errorNotif = newErrorNotification(rv.notificationArea)
	rv.togglePassword()
	rv.window.SetTitle(i18n.Localf("Room: [%s]", rv.roomJid.String()))
}

func (rv *mucRoomView) togglePassword() {
	doInUIThread(func() {
		value := rv.chkPassword.GetActive()
		rv.labelPassword.SetSensitive(value)
		rv.textPassword.SetSensitive(value)
	})
}

func (rv *mucRoomView) hasValidNickname() bool {
	nickName, _ := rv.textNickname.GetText()
	return len(nickName) > 0
}

func (rv *mucRoomView) hasValidPassword() bool {
	value := rv.chkPassword.GetActive()
	if !value {
		return true
	}
	password, _ := rv.textPassword.GetText()
	return len(password) > 0
}

func (rv *mucRoomView) validateInput() {
	sensitiveValue := rv.hasValidNickname() && rv.hasValidPassword()
	rv.btnAcceptJoin.SetSensitive(sensitiveValue)
}

func (rv *mucRoomView) togglePanelView() {
	doInUIThread(func() {
		value := rv.boxJoinRoomView.IsVisible()
		rv.boxJoinRoomView.SetVisible(!value)
		rv.boxRoomView.SetVisible(value)
	})
}

func (rv *mucRoomView) onBtnJoinClicked() {
	nickName, _ := rv.textNickname.GetText()
	go rv.account.session.JoinRoom(rv.roomJid, nickName)
	rv.togglePanelView()
}

func (u *gtkUI) mucShowRoom(a *account, rjid jid.Bare) {
	view := &mucRoomView{}

	view.account = a
	view.roomJid = rjid
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
