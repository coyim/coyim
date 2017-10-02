package gtki

import "github.com/coyim/gotk3adapter/glibi"

type Application interface {
	glibi.Application

	GetActiveWindow() Window
}

func AssertApplication(_ Application) {}
