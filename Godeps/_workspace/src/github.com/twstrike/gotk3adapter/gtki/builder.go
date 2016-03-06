package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type Builder interface {
	glibi.Object

	AddFromString(string) error
	ConnectSignals(map[string]interface{})
	GetObject(string) (glibi.Object, error)
}

func AssertBuilder(_ Builder) {}
