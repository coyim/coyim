package gui

import (
	"sync"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomView struct {
	builder *builder

	id         string
	generation int
	updateLock sync.RWMutex

	window gtki.Window `gtk-widget:"MUCRoom"`
	//room   *room

	boxJoinRoomView  gtki.Box         `gtk-widget:"boxJoinRoomView"`
	labelNickname    gtki.Label       `gtk-widget:"labelNickname"`
	txtNickname      gtki.Entry       `gtk-widget:"txtNickname"`
	chkPassword      gtki.CheckButton `gtk-widget:"chkPassword"`
	labelPassword    gtki.Label       `gtk-widget:"labelPassword"`
	txtPassword      gtki.Entry       `gtk-widget:"txtPassword"`
	notificationArea gtki.Box         `gtk-widget:"boxNotificationArea"`
	notification     gtki.InfoBar
	errorNotif       *errorNotification

	boxRoomView gtki.Box    `gtk-widget:"boxRoomView"`
	panel       gtki.Box    `gtk-widget:"panel"`
	panelToggle gtki.Button `gtk-widget:"panel-toggle"`

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
	rv.builder = newBuilder("MUCRoom")
	panicOnDevError(rv.builder.bindObjects(rv))
	rv.errorNotif = newErrorNotification(rv.notificationArea)
	rv.togglePassword()
}

func (rv *mucRoomView) togglePassword() {
	doInUIThread(func() {
		value := rv.chkPassword.GetActive()
		rv.labelPassword.SetSensitive(!value)
		rv.txtPassword.SetSensitive(!value)
	})
}

//
// Custom GTK Events
//

//
func (rv *mucRoomView) onShowWindow() {
	//TODO: add necessary calls here
}

func (rv *mucRoomView) onCloseWindow() {
	//TODO: add necessary calls here
}

func (rv *mucRoomView) onPasswordChecked() {
	rv.togglePassword()
}

func (rv *mucRoomView) onBtnJoinClicked() {
	//TODO: calls the XMPP logic to join a room here
}

// mucShowRoom should be called from the UI thread
func (u *gtkUI) mucShowRoom(userAccount *account) {
	view := &mucRoomView{}

	view.userAccount = userAccount
	view.init()

	view.builder.ConnectSignals(map[string]interface{}{
		"on_close_window_signal": func() {},
		"on_show_window_signal": func() {
			view.onShowWindow()
		},
		"on_chk_password_checked_signal": func() {
			view.onPasswordChecked()
		},
		"on_btn_cancel_clicked_signal": view.window.Destroy,
		"on_btn_join_clicked_signal": func() {
			view.onBtnJoinClicked()
		},
	})

	u.connectShortcutsChildWindow(view.window)

	view.window.SetTransientFor(u.window)
	view.window.Show()
}
