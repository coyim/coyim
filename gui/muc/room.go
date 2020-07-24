package muc

import "github.com/coyim/gotk3adapter/gtki"

type mucRoomUI struct {
	panel       gtki.Box    `gtk-widget:"panel"`
	panelToggle gtki.Button `gtk-widget:"panel-toggle"`

	room gtki.Box `gtk-widget:"room"`

	roomPanelOpen  bool
	roomViewActive bool
}

func (m *mucRoomUI) togglePanel() {
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

func (m *mucRoomUI) closeRoomWindow() {
	if !m.roomViewActive {
		return
	}

	m.roomViewActive = false
}

func (m *mucUI) openRoomView(id string) {
	// TODO: show the room window
	panic("Show room window")
}
