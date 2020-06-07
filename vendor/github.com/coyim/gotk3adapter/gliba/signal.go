package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type signal struct {
	*glib.Signal
}

func WrapSignalSimple(s *glib.Signal) glibi.Signal {
	if s == nil {
		return nil
	}
	return &signal{s}
}

func wrapSignal(s *glib.Signal, e error) (glibi.Signal, error) {
	return WrapSignalSimple(s), nil
}

func UnwrapSignal(v glibi.Signal) *glib.Signal {
	if v == nil {
		return nil
	}
	return v.(*signal).Signal
}
