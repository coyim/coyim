// +build !gtk_3_6,!gtk_3_8,!gtk_3_10

package gtka

import "github.com/coyim/gotk3adapter/gtki"

func (v *box) SetCenterWidget(v2 gtki.Widget) {
	v.internal.SetCenterWidget(UnwrapWidget(v2))
}

func (v *box) GetCenterWidget() gtki.Widget {
	vx := v.internal.GetCenterWidget()
	return Wrap(vx).(gtki.Widget)
}
