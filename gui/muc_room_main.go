package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewMain struct {
	ident jid.Bare
	ac    *account
	log   coylog.Logger

	content gtki.Box `gtk-widget:"boxRoomView"`
	topBox  gtki.Box `gtk-widget:"roomViewTop"`
	roomBox gtki.Box `gtk-widget:"room"`
	paneBox gtki.Box `gtk-widget:"panel"`

	main   gtki.Box
	panel  gtki.Box
	top    gtki.Box
	parent gtki.Box
}

func (v *roomView) initRoomMain() {
	v.main = newRoomMainView(
		v.account,
		v.identity,
		v.conv.view,
		v.roster.view,
		v.toolbar.view,
		v.content,
	)
}

func newRoomMainView(a *account, rid jid.Bare, main, panel, top, parent gtki.Box) *roomViewMain {
	m := &roomViewMain{
		ident:  rid,
		ac:     a,
		main:   main,
		panel:  panel,
		top:    top,
		parent: parent,
	}

	builder := newBuilder("MUCRoomMain")
	panicOnDevError(builder.bindObjects(m))

	m.log = a.log.WithField("room", m.ident)

	m.roomBox.SetHExpand(true)
	m.content.SetHExpand(true)

	m.roomBox.Add(m.main)
	m.paneBox.Add(m.panel)
	m.topBox.Add(m.top)

	m.parent.Add(m.content)

	return m
}

func (v *roomViewMain) show() {
	v.content.Show()
}

func (v *roomViewMain) hide() {}
