package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type spinner struct {
	*widget
	internal *gtk.Spinner
}

func wrapSpinnerSimple(v *gtk.Spinner) *spinner {
	if v == nil {
		return nil
	}

	return &spinner{wrapWidgetSimple(&v.Widget), v}
}

func wrapSpinner(v *gtk.Spinner, e error) (*spinner, error) {
	return wrapSpinnerSimple(v), e
}

func unwrapSpinner(v gtki.Spinner) *gtk.Spinner {
	if v == nil {
		return nil
	}
	return v.(*spinner).internal
}

func (v *spinner) Start() {
	v.internal.Start()
}

func (v *spinner) Stop() {
	v.internal.Stop()
}
