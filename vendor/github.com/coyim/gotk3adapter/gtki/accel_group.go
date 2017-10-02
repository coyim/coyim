package gtki

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/glibi"
)

type AccelGroup interface {
	glibi.Object

	Connect2(uint, gdki.ModifierType, AccelFlags, interface{})
}

func AssertAccelGroup(_ AccelGroup) {}
