package gtki

import "github.com/twstrike/gotk3adapter/gdki"

type Menu interface {
	MenuShell

	PopupAtMouseCursor(Menu, MenuItem, int, uint32)
	PopupAtPointer(gdki.Event)
}

func AssertMenu(_ Menu) {}
