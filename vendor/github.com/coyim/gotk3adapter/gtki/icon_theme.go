package gtki

import "github.com/coyim/gotk3adapter/glibi"

type IconTheme interface {
	glibi.Object

	AddResourcePath(string)
	AppendSearchPath(string)
	GetExampleIconName() string
	HasIcon(string) bool
	PrependSearchPath(string)
}

func AssertIconTheme(_ IconTheme) {}
