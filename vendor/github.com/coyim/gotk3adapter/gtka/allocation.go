package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type allocation struct {
	internal *gtk.Allocation
}

func WrapAllocationSimple(v *gtk.Allocation) gtki.Allocation {
	if v == nil {
		return nil
	}
	return &allocation{v}
}

func WrapAllocation(v *gtk.Allocation, e error) (gtki.Allocation, error) {
	return WrapAllocationSimple(v), e
}

func UnwrapAllocation(v gtki.Allocation) *gtk.Allocation {
	if v == nil {
		return nil
	}
	return v.(*allocation).internal
}

func (v *allocation) GetY() int {
	return v.internal.GetY()
}
