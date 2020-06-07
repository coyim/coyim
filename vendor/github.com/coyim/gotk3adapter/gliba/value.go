package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type value struct {
	*glib.Value
}

func WrapValueSimple(v *glib.Value) glibi.Value {
	if v == nil {
		return nil
	}
	return &value{v}
}

func WrapValue(v *glib.Value, e error) (glibi.Value, error) {
	return WrapValueSimple(v), e
}

func UnwrapValue(v glibi.Value) *glib.Value {
	if v == nil {
		return nil
	}
	return v.(*value).Value
}
