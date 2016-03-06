package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type Application interface {
	glibi.Application

	GetActiveWindow() Window
}

func AssertApplication(_ Application) {}
