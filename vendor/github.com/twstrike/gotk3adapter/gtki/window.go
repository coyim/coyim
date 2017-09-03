package gtki

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/twstrike/gotk3adapter/gdki"
)

type Window interface {
	Bin

	AddAccelGroup(AccelGroup)
	GetTitle() string
	GetSize() (int, int)
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
	AddMnemonic(uint, Widget)
	RemoveMnemonic(uint, Widget)
	ActivateMnemonic(uint, gdki.ModifierType) bool
	GetMnemonicModifier() gdk.ModifierType
	SetMnemonicModifier(gdki.ModifierType)
}

func AssertWindow(_ Window) {}
