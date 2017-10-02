package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type adjustment struct {
	*gliba.Object
	internal *gtk.Adjustment
}

func wrapAdjustmentSimple(v *gtk.Adjustment) *adjustment {
	if v == nil {
		return nil
	}
	return &adjustment{gliba.WrapObjectSimple(v.Object), v}
}

func wrapAdjustment(v *gtk.Adjustment, e error) (*adjustment, error) {
	return wrapAdjustmentSimple(v), e
}

func unwrapAdjustment(v gtki.Adjustment) *gtk.Adjustment {
	if v == nil {
		return nil
	}
	return v.(*adjustment).internal
}

func (v *adjustment) GetLower() float64 {
	return v.internal.GetLower()
}

func (v *adjustment) GetPageSize() float64 {
	return v.internal.GetPageSize()
}

func (v *adjustment) GetUpper() float64 {
	return v.internal.GetUpper()
}

func (v *adjustment) SetValue(v1 float64) {
	v.internal.SetValue(v1)
}
