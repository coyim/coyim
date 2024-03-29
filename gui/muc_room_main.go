package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewMain struct {
	main  gtki.Box
	panel gtki.Box
	top   gtki.Box

	content gtki.Box `gtk-widget:"room-view-box"`
	topBox  gtki.Box `gtk-widget:"room-view-top"`
	roomBox gtki.Box `gtk-widget:"room-view-content"`
	paneBox gtki.Box `gtk-widget:"room-view-panel"`
}

func (v *roomView) newRoomMainView() *roomViewMain {
	m := &roomViewMain{
		main:  v.conv.view,
		panel: v.roster.view,
		top:   v.toolbar.view,
	}

	m.initBuilder()
	m.initDefaults()

	return m
}

func (m *roomViewMain) initBuilder() {
	buildUserInterface("MUCRoomMain", m)
}

func (m *roomViewMain) initDefaults() {
	m.roomBox.Add(m.main)
	m.paneBox.Add(m.panel)
	m.topBox.Add(m.top)
}
