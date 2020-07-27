package muc

import (
	"errors"

	"github.com/coyim/gotk3adapter/gtki"
)

type roomUI struct {
	id string

	panel       gtki.Box    `gtk-widget:"panel"`
	panelToggle gtki.Button `gtk-widget:"panel-toggle"`

	room gtki.Box `gtk-widget:"room"`

	roomPanelOpen  bool
	roomViewActive bool

	window  gtki.Window
	builder *builder
}

func (u *gtkUI) openRoomView(id string) {
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

	u.doInUIThread(win.Show)

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
