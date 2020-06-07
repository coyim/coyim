package gtki

type ApplicationWindow interface {
	Window

	SetShowMenubar(bool)
	GetShowMenubar() bool
	GetID() uint
}

func AssertApplicationWindow(_ ApplicationWindow) {}
