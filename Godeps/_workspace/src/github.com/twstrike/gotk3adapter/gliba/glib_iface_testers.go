package gliba

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

func init() {
	glibi.AssertGlib(&RealGlib{})
	glibi.AssertApplication(&Application{})
	glibi.AssertObject(&Object{})
	glibi.AssertSignal(&signal{})
	glibi.AssertValue(&value{})
}
