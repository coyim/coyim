package gtki

import "github.com/coyim/gotk3adapter/gdki"

type Menu interface {
	MenuShell

	PopupAtPointer(gdki.Event)
}

func AssertMenu(_ Menu) {}
