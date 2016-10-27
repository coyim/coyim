package gtki

type ApplicationWindow interface {
	Window
}

func AssertApplicationWindow(_ ApplicationWindow) {}
