package gtki

import "github.com/coyim/gotk3adapter/glibi"

type CssProvider interface {
	glibi.Object

	LoadFromData(string) error
}

func AssertCssProvider(_ CssProvider) {}
