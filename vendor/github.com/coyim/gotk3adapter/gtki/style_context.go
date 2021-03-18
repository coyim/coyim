package gtki

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
)

type StyleContext interface {
	glibi.Object

	AddClass(string)
	RemoveClass(string)
	AddProvider(StyleProvider, uint)
	GetScreen() (gdki.Screen, error)
	GetProperty2(string, StateFlags) (interface{}, error)
}

func AssertStyleContext(_ StyleContext) {}
