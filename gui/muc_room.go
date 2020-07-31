package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/session/muc"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomView struct {
	builder *builder

	id         string
	generation int
	updateLock sync.RWMutex

	window gtki.Window `gtk-widget:"MUCRoomWindow"`
	//room   *room

	boxJoinRoomView  gtki.Box         `gtk-widget:"boxJoinRoomView"`
	txtNickname      gtki.Entry       `gtk-widget:"txtNickname"`
	chkPassword      gtki.CheckButton `gtk-widget:"chkPassword"`
	labelPassword    gtki.Label       `gtk-widget:"labelPassword"`
	txtPassword      gtki.Entry       `gtk-widget:"txtPassword"`
	btnAcceptJoin    gtki.Button      `gtk-widget:"btnAcceptJoin"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification

	boxRoomView gtki.Box `gtk-widget:"boxRoomView"`
	//panel       gtki.Box    `gtk-widget:"panel"`
	//panelToggle gtki.Button `gtk-widget:"panel-toggle"`

	/*
		roomPanelOpen    bool
		roomViewActive   bool
		roomConversation gtki.Box `gtk-widget:"room"`

		windowTitle       gtki.Label `gtk-widget:"window-title"`
		windowDescription gtki.Label `gtk-widget:"window-description"`

		roomMembers      gtki.ScrolledWindow `gtk-widget:"room-members"`
		roomMembersModel gtki.ListStore      `gtk-widget:"room-members-model"`
		roomMembersView  gtki.TreeView       `gtk-widget:"room-members-tree"`
	*/
	// using the room jid for a moment, this should be an interface with all the necessary room information
	roomInfo    *muc.RoomListing
	userAccount *account
}

/*
type roomRole string

const (
	roleAdministrator roomRole = "administrator"
	roleModerator     roomRole = "moderator"
	roleNone          roomRole = "none"
	roleParticipant   roomRole = "participant"
	roleVisitor       roomRole = "visitor"
)

type room struct {
	*rosterItem
	description string
	members     membersList
}
*/

/*
func (u *gtkUI) openRoomView(id string) {
	room, err := u.roomsServer.byID(id)
	if err != nil {
		panic(err.Error())
	}

	r2, err := u.roomWindowByID(id)
	if err == nil {
		r2.window.Present()
		return
	}

	builder := newBuilder("room")

	win := builder.get("roomWindow").(gtki.Window)

	win.SetTitle(id)

	r := &roomUI{
		id:      id,
		builder: builder,
		window:  win,
		room:    room,
		u:       u,
	}

	panicOnDevError(builder.bindObjects(r))

	builder.ConnectSignals(map[string]interface{}{
		"on_close": func() {
			u.removeRoomWindow(id)
		},
		"on_conversation_close": func() {
			u.doInUIThread(win.Destroy)
		},
		"on_toggle_panel": r.togglePanel,
	})

	u.doInUIThread(func() {
		win.Show()
		r.windowTitle.SetText(r.room.displayName())
		r.windowDescription.SetText(r.room.displayDescription())
	})

	// TODO: improve this if required for future interactions
	// in this mockup
	r.showRoomMembers()

	u.addNewRoomWindow(id, r)
}
*/

/*
func (u *roomUI) togglePanel() {
	isOpen := !u.roomPanelOpen

	var toggleLabel string
	if isOpen {
		toggleLabel = "Hide panel"
	} else {
		toggleLabel = "Show panel"
	}
	u.panelToggle.SetProperty("label", toggleLabel)
	u.panel.SetVisible(isOpen)
	u.roomPanelOpen = isOpen
}

func (u *roomUI) closeRoomWindow() {
	if !u.roomViewActive {
		return
	}

	u.roomViewActive = false
}

func (u *gtkUI) addNewRoomWindow(id string, r *roomUI) {
	_, err := u.roomWindowByID(id)
	if err != nil {
		u.roomUI[id] = r
	}
}

func (u *gtkUI) removeRoomWindow(id string) {
	if _, ok := u.roomUI[id]; ok {
		delete(u.roomUI, id)
	}
}

func (u *gtkUI) roomWindowByID(id string) (*roomUI, error) {
	if r, ok := u.roomUI[id]; ok {
		return r, nil
	}
	return nil, errors.New("room window don't exists")
}

func (r *room) displayDescription() string {
	return r.description
}
*/

// init: Initilize the rooms view
func (rv *mucRoomView) init() {
	rv.builder = newBuilder("MUCRoomWindow")
	panicOnDevError(rv.builder.bindObjects(rv))
	rv.errorNotif = newErrorNotification(rv.notificationArea)
	rv.togglePassword()

	doInUIThread(func() {
		rv.window.SetTitle(fmt.Sprintf("Room: [%s]", rv.roomInfo.Jid.String()))
	})
}

// tooglePassword activate/deactivate the password fields
func (rv *mucRoomView) togglePassword() {
	doInUIThread(func() {
		value := rv.chkPassword.GetActive()
		rv.labelPassword.SetSensitive(value)
		rv.txtPassword.SetSensitive(value)
	})
}

func (rv *mucRoomView) validateInput() {
	doInUIThread(func() {
		value := rv.chkPassword.GetActive()
		nickName, _ := rv.txtNickname.GetText()
		password, _ := rv.txtPassword.GetText()
		sensitiveValue := len(nickName) > 0 && ((len(password) > 0 && value) || !value)
		rv.btnAcceptJoin.SetSensitive(sensitiveValue)
	})
}

//
// Custom GTK Events
//

// onShowWindow apply some actions when the view is showed
func (rv *mucRoomView) onShowWindow() {
	//TODO: add necessary calls here
}

// onCloseWindow apply some actions when the view is hidden
func (rv *mucRoomView) onCloseWindow() {
	//TODO: add necessary calls here
}

// onBtnJoinClicked event handler for the click event on the button join
func (rv *mucRoomView) onBtnJoinClicked() {
	//TODO: calls the XMPP logic to join a room here
	//nickName, _ := rv.txtNickname.GetText()
	//rv.userAccount.session.JoinRoom(rv.roomJid, nickName)
}

// mucShowRoom should be called from the UI thread
func (u *gtkUI) mucShowRoom(userAccount *account, rl *muc.RoomListing) {
	view := &mucRoomView{}

	view.userAccount = userAccount
	view.roomInfo = rl
	view.init()

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal": func() {},
		"on_show_window_signal": func() {
			view.validateInput()
		},
		"on_txt_nickname_changed_signal": func() {
			view.validateInput()
		},
		"on_txt_password_changed_signal": func() {
			view.validateInput()
		},
		"on_chk_password_checked_signal": func() {
			view.togglePassword()
			view.validateInput()
		},
		"on_btn_cancel_join_clicked_signal": view.window.Destroy,
		"on_btn_accept_join_clicked_signal": func() {
			view.onBtnJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.window)

	view.window.SetTransientFor(u.window)
	view.window.Show()
}
