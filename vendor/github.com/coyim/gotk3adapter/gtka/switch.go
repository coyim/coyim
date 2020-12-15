package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

// We can't use `switch` here because it's a reserved word
type zwitch struct {
	*widget
	internal *gtk.Switch
}

func WrapSwitchSimple(v *gtk.Switch) gtki.Switch {
	if v == nil {
		return nil
	}
	return &zwitch{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapSwitch(v *gtk.Switch, e error) (gtki.Switch, error) {
	return WrapSwitchSimple(v), e
}

func UnwrapSwitch(v gtki.Switch) *gtk.Switch {
	if v == nil {
		return nil
	}
	return v.(*zwitch).internal
}

func (v *zwitch) SetActive(v1 bool) {
	v.internal.SetActive(v1)
}

func (v *zwitch) GetActive() bool {
	return v.internal.GetActive()
}
