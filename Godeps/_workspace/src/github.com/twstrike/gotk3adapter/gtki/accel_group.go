package gtki

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
)

type AccelGroup interface {
	glibi.Object

	Connect2(uint, gdki.ModifierType, AccelFlags, interface{})
}

func AssertAccelGroup(_ AccelGroup) {}
