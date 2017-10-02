package gliba

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/coyim/gotk3adapter/glibi"
)

type value struct {
	*glib.Value
}

func wrapValueSimple(v *glib.Value) *value {
	if v == nil {
		return nil
	}
	return &value{v}
}

func WrapValue(v *glib.Value, e error) (*value, error) {
	return wrapValueSimple(v), e
}

func unwrapValue(v glibi.Value) *glib.Value {
	if v == nil {
		return nil
	}
	return v.(*value).Value
}
