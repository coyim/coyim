package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"

type Window interface {
	Bin

	AddAccelGroup(AccelGroup)
	GetTitle() string
	HasToplevelFocus() bool
	IsActive() bool
	Fullscreen()
	Unfullscreen()
	Present()
	Resize(int, int)
	SetApplication(Application)
	SetIcon(gdki.Pixbuf)
	SetTitle(string)
	SetTitlebar(Widget) // Since 3.10
	SetTransientFor(Window)
	SetUrgencyHint(bool)
}

func AssertWindow(_ Window) {}
