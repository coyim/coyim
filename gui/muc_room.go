package gui

import (
	"fmt"
	"sync"

	"github.com/coyim/coyim/xmpp/jid"

	"github.com/coyim/coyim/session/muc"

	"github.com/coyim/gotk3adapter/gtki"
)

type mucRoomView struct {
	builder *builder

	id         string
	generation int
	updateLock sync.RWMutex

	window gtki.Window `gtk-widget:"room-window"`
	//room   *room

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
	roomJid  jid.Bare
	roomInfo *muc.RoomListing
	account  *account
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
		rv.window.SetTitle(fmt.Sprintf("Room: [%s]", rv.roomJid.String()))
	})
}

// tooglePassword activate/deactivate the password fields
func (rv *mucRoomView) togglePassword() {
	doInUIThread(func() {
		value := rv.chkPassword.GetActive()
		rv.labelPassword.SetSensitive(value)
		rv.textPassword.SetSensitive(value)
	})
}

// hasValidNickname checking if the nickname has entered
func (rv *mucRoomView) hasValidNickname() bool {
	nickName, _ := rv.textNickname.GetText()
	return len(nickName) > 0
}

// hasValidPassword checking if the password has checked and entered
func (rv *mucRoomView) hasValidPassword() bool {
	value := rv.chkPassword.GetActive()
	if !value {
		return true
	}
	password, _ := rv.textPassword.GetText()
	return len(password) > 0
}

// validateInput checking if the button join must be enable in order to execute the action
func (rv *mucRoomView) validateInput() {
	doInUIThread(func() {
		sensitiveValue := rv.hasValidNickname() && rv.hasValidPassword()
		rv.btnAcceptJoin.SetSensitive(sensitiveValue)
	})
}

// togglePanelView toggle the view between the join panel and chat panel
func (rv *mucRoomView) togglePanelView() {
	doInUIThread(func() {
		value := rv.boxJoinRoomView.IsVisible()
		rv.boxJoinRoomView.SetVisible(!value)
		rv.boxRoomView.SetVisible(value)
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
	defer func() {
		rv.togglePanelView()
	}()
	nickName, _ := rv.textNickname.GetText()
	go rv.account.session.JoinRoom(rv.roomJid, nickName)
}

// mucShowRoom should be called from the UI thread
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

	view.window.SetTransientFor(u.window)
	view.window.Show()
}
