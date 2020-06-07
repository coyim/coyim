package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type progressBar struct {
	*widget
	internal *gtk.ProgressBar
}

func WrapProgressBarSimple(v *gtk.ProgressBar) gtki.ProgressBar {
	if v == nil {
		return nil
	}

	return &progressBar{WrapWidgetSimple(&v.Widget).(*widget), v}
}

func WrapProgressBar(v *gtk.ProgressBar, e error) (gtki.ProgressBar, error) {
	return WrapProgressBarSimple(v), e
}

func UnwrapProgressBar(v gtki.ProgressBar) *gtk.ProgressBar {
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
