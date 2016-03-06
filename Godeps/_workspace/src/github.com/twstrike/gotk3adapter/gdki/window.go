package gdki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type Window interface {
	glibi.Object

	GetDesktop() uint32
	MoveToDesktop(uint32)
}

func AssertWindow(_ Window) {}
