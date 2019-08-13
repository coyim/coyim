package gtki

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type Window interface {
	Bin

	ActivateMnemonic(uint, gdki.ModifierType) bool
	AddAccelGroup(AccelGroup)
	AddMnemonic(uint, Widget)
	Deiconify()
	Fullscreen()
	GetMnemonicModifier() gdk.ModifierType
	GetTitle() string
	GetSize() (int, int)
	HasToplevelFocus() bool
	Iconify()
	IsActive() bool
	Maximize()
	Present()
	Resize(int, int)
	RemoveMnemonic(uint, Widget)
	SetApplication(Application)
	SetDecorated(bool)
	SetIcon(gdki.Pixbuf)
	SetMnemonicModifier(gdki.ModifierType)
	SetTitle(string)
	SetTitlebar(Widget) // Since 3.10
	SetTransientFor(Window)
	GetTransientFor() (Window, error)
	SetUrgencyHint(bool)
	Unfullscreen()
	Unmaximize()
	GetPosition() (int, int)
	Move(int, int)
}

func AssertWindow(_ Window) {}
