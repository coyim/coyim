package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type StyleContext interface {
	glibi.Object

	AddClass(string)
	AddProvider(StyleProvider, uint)
	GetProperty2(string, StateFlags) (interface{}, error)
}

func AssertStyleContext(_ StyleContext) {}
