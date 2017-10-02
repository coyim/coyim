package gtki

import "github.com/coyim/gotk3adapter/glibi"

type Builder interface {
	glibi.Object

	AddFromResource(string) error
	AddFromString(string) error
	ConnectSignals(map[string]interface{})
	GetObject(string) (glibi.Object, error)
}

func AssertBuilder(_ Builder) {}
