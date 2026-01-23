package gui

import "github.com/coyim/gotk3adapter/glibi"

func (u *gtkUI) addAction(m glibi.ActionMap, name string, f func()) {
	act := g.glib.SimpleActionNew(name, nil)
	ignore(act.Connect("activate", f))
	m.AddAction(act)
}
