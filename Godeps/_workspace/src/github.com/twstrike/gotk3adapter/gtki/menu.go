package gtki

type Menu interface {
	MenuShell

	PopupAtMouseCursor(Menu, MenuItem, int, uint32)
}

func AssertMenu(_ Menu) {}
