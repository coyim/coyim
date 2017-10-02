package gtki

import "github.com/coyim/gotk3adapter/glibi"

type Settings interface {
	glibi.Object
}

func AssertSettings(_ Settings) {}
