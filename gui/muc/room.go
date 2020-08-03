package muc

import (
	"errors"

	"github.com/coyim/gotk3adapter/gtki"
)

type roomUI struct {
	id      string
	u       *gtkUI
	builder *builder
	window  gtki.Window
	room    *room

	panel       gtki.Box    `gtk-widget:"panel"`
	panelToggle gtki.Button `gtk-widget:"panel-toggle"`

	roomPanelOpen    bool
	roomViewActive   bool
	roomConversation gtki.Box `gtk-widget:"room"`

	windowTitle       gtki.Label `gtk-widget:"window-title"`
	windowDescription gtki.Label `gtk-widget:"window-description"`

	roomMembers      gtki.ScrolledWindow `gtk-widget:"room-members"`
	roomMembersModel gtki.ListStore      `gtk-widget:"room-members-model"`
	roomMembersView  gtki.TreeView       `gtk-widget:"room-members-tree"`
}

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

func (u *roomUI) togglePanel() {
	isOpen := !u.roomPanelOpen

	var toggleLabel string
	if isOpen {
		toggleLabel = "Hide panel"
	} else {
		toggleLabel = "Show panel"
	}
	_ = u.panelToggle.SetProperty("label", toggleLabel)
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
