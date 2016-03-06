package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"

type Window interface {
	Bin

	AddAccelGroup(AccelGroup)
	IsActive() bool
	Present()
	Resize(int, int)
	SetApplication(Application)
	SetIcon(gdki.Pixbuf)
	SetTitle(string)
	SetTitlebar(Widget) // Since 3.10
	SetTransientFor(Window)
}

func AssertWindow(_ Window) {}
