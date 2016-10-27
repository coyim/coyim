package gtki

import "github.com/twstrike/gotk3adapter/glibi"

type Settings interface {
	glibi.Object
}

func AssertSettings(_ Settings) {}
