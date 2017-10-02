package pangoa

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/gotk3/gotk3/pango"
)

func init() {
	gliba.AddWrapper(WrapLocal)

	gliba.AddUnwrapper(UnwrapLocal)
}

func WrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *pango.FontDescription:
		val := wrapFontDescriptionSimple(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}

func UnwrapLocal(o interface{}) (interface{}, bool) {
	switch oo := o.(type) {
	case *fontDescription:
		val := unwrapFontDescription(oo)
		if val == nil {
			return nil, true
		}
		return val, true
	default:
		return nil, false
	}
}
