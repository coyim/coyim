package gdki

import "github.com/coyim/gotk3adapter/glibi"

type Window interface {
	glibi.Object

	GetDesktop() uint32
	MoveToDesktop(uint32)
}

func AssertWindow(_ Window) {}
