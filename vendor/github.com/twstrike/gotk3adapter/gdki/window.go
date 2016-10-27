package gdki

import "github.com/twstrike/gotk3adapter/glibi"

type Window interface {
	glibi.Object

	GetDesktop() uint32
	MoveToDesktop(uint32)
}

func AssertWindow(_ Window) {}
