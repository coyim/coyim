package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type progressBar struct {
	*widget
	internal *gtk.ProgressBar
}

func wrapProgressBarSimple(v *gtk.ProgressBar) *progressBar {
	if v == nil {
		return nil
	}

	return &progressBar{wrapWidgetSimple(&v.Widget), v}
}

func wrapProgressBar(v *gtk.ProgressBar, e error) (*progressBar, error) {
	return wrapProgressBarSimple(v), e
}

func unwrapProgressBar(v gtki.ProgressBar) *gtk.ProgressBar {
	if v == nil {
		return nil
	}
	return v.(*progressBar).internal
}

func (v *progressBar) SetFraction(v1 float64) {
	v.internal.SetFraction(v1)
}

func (v *progressBar) GetFraction() float64 {
	return v.internal.GetFraction()
}

func (v *progressBar) SetShowText(v1 bool) {
	v.internal.SetShowText(v1)
}

func (v *progressBar) GetShowText() bool {
	return v.internal.GetShowText()
}

func (v *progressBar) SetText(v1 string) {
	v.internal.SetText(v1)
}
