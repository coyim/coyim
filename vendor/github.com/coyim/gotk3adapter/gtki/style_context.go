package gtki

import "github.com/coyim/gotk3adapter/glibi"

type StyleContext interface {
	glibi.Object

	AddClass(string)
	AddProvider(StyleProvider, uint)
	GetProperty2(string, StateFlags) (interface{}, error)
}

func AssertStyleContext(_ StyleContext) {}
