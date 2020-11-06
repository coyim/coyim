package gtki

type MenuButton interface {
	Bin

	SetPopover(Popover)
}

func AssertMenuButton(_ MenuButton) {}
