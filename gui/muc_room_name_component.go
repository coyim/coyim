package gui

import (
	"github.com/coyim/coyim/xmpp/jid"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomNameComponent struct {
	entry gtki.Entry
}

func (u *gtkUI) createRoomNameComponent(e gtki.Entry) *roomNameComponent {
	return &roomNameComponent{e}
}

func (rn *roomNameComponent) currentName() jid.Local {
	return jid.NewLocal(rn.currentNameValue())
}

func (rn *roomNameComponent) currentNameValue() string {
	t, _ := rn.entry.GetText()
	return t
}

func (rn *roomNameComponent) setCurrentNameValue(t string) {
	rn.entry.SetText(t)
}

func (rn *roomNameComponent) reset() {
	rn.setCurrentNameValue("")
}

func (rn *roomNameComponent) hasNameValue() bool {
	return rn.currentNameValue() != ""
}

func (rn *roomNameComponent) disableNameInput() {
	rn.entry.SetSensitive(false)
}

func (rn *roomNameComponent) enableNameInput() {
	rn.entry.SetSensitive(true)
}
