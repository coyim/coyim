package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type signal struct {
	*glib.Signal
}

func wrapSignalSimple(s *glib.Signal) *signal {
	if s == nil {
		return nil
	}
	return &signal{s}
}

func wrapSignal(s *glib.Signal, e error) (*signal, error) {
	return wrapSignalSimple(s), nil
}

func unwrapSignal(v glibi.Signal) *glib.Signal {
	if v == nil {
		return nil
	}
	return v.(*signal).Signal
}
