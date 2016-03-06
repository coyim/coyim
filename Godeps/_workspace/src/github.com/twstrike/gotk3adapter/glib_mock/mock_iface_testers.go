package glib_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

func init() {
	glibi.AssertGlib(&Mock{})
	glibi.AssertApplication(&MockApplication{})
	glibi.AssertObject(&MockObject{})
	glibi.AssertSignal(&MockSignal{})
	glibi.AssertValue(&MockValue{})
}
