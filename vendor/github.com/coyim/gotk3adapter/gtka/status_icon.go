package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
)

type statusIcon struct {
	*gliba.Object
	internal *gotk3extra.StatusIcon
}

func WrapStatusIconSimple(v *gotk3extra.StatusIcon) gtki.StatusIcon {
	if v == nil {
		return nil
	}
	return &statusIcon{gliba.WrapObjectSimple(v.Object), v}
}

func WrapStatusIcon(v *gotk3extra.StatusIcon, e error) (gtki.StatusIcon, error) {
	return WrapStatusIconSimple(v), e
}

func UnwrapStatusIcon(v gtki.StatusIcon) *gotk3extra.StatusIcon {
	if v == nil {
		return nil
	}
	return v.(*statusIcon).internal
}

func (v *statusIcon) SetFromFile(filename string) {
	v.internal.SetFromFile(filename)
}

func (v *statusIcon) SetFromIconName(iconName string) {
	v.internal.SetFromIconName(iconName)
}

func (v *statusIcon) SetFromPixbuf(pixbuf gdki.Pixbuf) {
	v.internal.SetFromPixbuf(gdka.UnwrapPixbuf(pixbuf))
}

func (v *statusIcon) SetTooltipText(text string) {
	v.internal.SetTooltipText(text)
}

func (v *statusIcon) GetTooltipText() string {
	return v.internal.GetTooltipText()
}

func (v *statusIcon) SetTooltipMarkup(markup string) {
	v.internal.SetTooltipMarkup(markup)
}

func (v *statusIcon) GetTooltipMarkup() string {
	return v.internal.GetTooltipMarkup()
}

func (v *statusIcon) SetHasTooltip(hasTooltip bool) {
	v.internal.SetHasTooltip(hasTooltip)
}

func (v *statusIcon) GetTitle() string {
	return v.internal.GetTitle()
}

func (v *statusIcon) SetName(name string) {
	v.internal.SetName(name)
}

func (v *statusIcon) SetVisible(visible bool) {
	v.internal.SetVisible(visible)
}

func (v *statusIcon) GetVisible() bool {
	return v.internal.GetVisible()
}

func (v *statusIcon) IsEmbedded() bool {
	return v.internal.IsEmbedded()
}

func (v *statusIcon) GetX11WindowID() uint32 {
	return v.internal.GetX11WindowID()
}

func (v *statusIcon) GetHasTooltip() bool {
	return v.internal.GetHasTooltip()
}

func (v *statusIcon) SetTitle(title string) {
	v.internal.SetTitle(title)
}

func (v *statusIcon) GetIconName() string {
	return v.internal.GetIconName()
}

func (v *statusIcon) GetSize() int {
	return v.internal.GetSize()
}
