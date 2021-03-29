package gtki

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
)

type Widget interface {
	glibi.Object

	Destroy()
	GetWindow() (gdki.Window, error)
	GrabFocus()
	GetAllocatedHeight() int
	GetAllocatedWidth() int
	GetName() (string, error)
	GetParent() (Widget, error)
	GetParentX() (Widget, error)
	GetStyleContext() (StyleContext, error)
	GrabDefault()
	SetCanFocus(bool)
	HasFocus() bool
	Hide()
	HideOnDelete()
	Map()
	SetHAlign(Align)
	SetHExpand(bool)
	SetVAlign(Align)
	SetVExpand(bool)
	SetMarginBottom(int)
	SetMarginTop(int)
	SetName(string)
	SetNoShowAll(bool)
	SetSensitive(bool)
	SetOpacity(float64)
	SetSizeRequest(int, int)
	SetTooltipText(string)
	SetVisible(bool)
	IsVisible() bool
	Show()
	ShowAll()
}

func AssertWidget(_ Widget) {}
