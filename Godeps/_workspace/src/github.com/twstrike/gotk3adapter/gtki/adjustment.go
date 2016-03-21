package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type Adjustment interface {
	glibi.Object

	GetLower() float64
	GetPageSize() float64
	GetUpper() float64
	SetValue(float64)
} // end of Adjustment

func AssertAdjustment(_ Adjustment) {}
