package gtki

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
)

type StatusIcon interface {
	glibi.Object

	SetFromFile(filename string)
	SetFromIconName(iconName string)
	SetFromPixbuf(pixbuf gdki.Pixbuf)
	SetTooltipText(text string)
	GetTooltipText() string
	SetTooltipMarkup(markup string)
	GetTooltipMarkup() string
	SetHasTooltip(hasTooltip bool)
	GetTitle() string
	SetName(name string)
	SetVisible(visible bool)
	GetVisible() bool
	IsEmbedded() bool
	GetX11WindowID() uint32
	GetHasTooltip() bool
	SetTitle(title string)
	GetIconName() string
	GetSize() int
}

func AssertStatusIcon(_ StatusIcon) {}
