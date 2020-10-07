package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewMain struct {
	main   gtki.Box
	panel  gtki.Box
	top    gtki.Box
	parent gtki.Box

	content gtki.Box `gtk-widget:"room-view-box"`
	topBox  gtki.Box `gtk-widget:"room-view-top"`
	roomBox gtki.Box `gtk-widget:"room-view-content"`
	paneBox gtki.Box `gtk-widget:"room-view-panel"`
}

func (v *roomView) initRoomMain() {
	v.main = newRoomMainView(
		v.conv.view,
		v.roster.view,
		v.toolbar.view,
		v.content,
	)
}

func newRoomMainView(main, panel, top, parent gtki.Box) *roomViewMain {
	m := &roomViewMain{
		main:   main,
		panel:  panel,
		top:    top,
		parent: parent,
	}

	m.initBuilder()
	m.initDefaults()

	return m
}

func (m *roomViewMain) initBuilder() {
	builder := newBuilder("MUCRoomMain")
	panicOnDevError(builder.bindObjects(m))
}

func (m *roomViewMain) initDefaults() {
	m.roomBox.Add(m.main)
	m.paneBox.Add(m.panel)
	m.topBox.Add(m.top)

	m.parent.Add(m.content)
}

func (m *roomViewMain) show() {
	m.content.Show()
}
