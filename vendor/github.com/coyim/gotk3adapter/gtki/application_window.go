package gtki

import "github.com/coyim/gotk3adapter/glibi"

type ApplicationWindow interface {
	Window
	glibi.ActionGroup
	glibi.ActionMap

	SetShowMenubar(bool)
	GetShowMenubar() bool
	GetID() uint
}

func AssertApplicationWindow(_ ApplicationWindow) {}
