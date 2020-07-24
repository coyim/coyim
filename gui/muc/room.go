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

func (m *mucUI) openRoomView(id string) {
	r2, err := m.roomWindowByID(id)
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
			m.removeRoomWindow(id)
		},
		"on_conversation_close": func() {
			m.doInUIThread(win.Destroy)
		},
		"on_toggle_panel": r.togglePanel,
	})

	m.doInUIThread(win.Show)

	m.addNewRoomWindow(id, r)
}

func (m *roomUI) togglePanel() {
	isOpen := !m.roomPanelOpen

	var toggleLabel string
	if isOpen {
		toggleLabel = "Hide panel"
	} else {
		toggleLabel = "Show panel"
	}
	m.panelToggle.SetProperty("label", toggleLabel)
	m.panel.SetVisible(isOpen)
	m.roomPanelOpen = isOpen
}

func (m *roomUI) closeRoomWindow() {
	if !m.roomViewActive {
		return
	}

	m.roomViewActive = false
}

func (m *mucUI) addNewRoomWindow(id string, r *roomUI) {
	_, err := m.roomWindowByID(id)
	if err != nil {
		m.roomWindows[id] = r
	}
}

func (m *mucUI) removeRoomWindow(id string) {
	if _, ok := m.roomWindows[id]; ok {
		delete(m.roomWindows, id)
	}
}

func (m *mucUI) roomWindowByID(id string) (*roomUI, error) {
	if r, ok := m.roomWindows[id]; ok {
		return r, nil
	}
	return nil, errors.New("room window don't exists")
}
