package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type spinner struct {
	*widget
	internal *gtk.Spinner
}

func WrapSpinnerSimple(v *gtk.Spinner) gtki.Spinner {
	if v == nil {
		return nil
	}

	return &spinner{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapSpinner(v *gtk.Spinner, e error) (gtki.Spinner, error) {
	return WrapSpinnerSimple(v), e
}

func UnwrapSpinner(v gtki.Spinner) *gtk.Spinner {
	if v == nil {
		return nil
	}
	return v.(*spinner).internal
}

func (*RealGtk) SpinnerNew() (gtki.Spinner, error) {
	return WrapSpinner(gtk.SpinnerNew())
}

func (v *spinner) Start() {
	v.internal.Start()
}

func (v *spinner) Stop() {
	v.internal.Stop()
}
