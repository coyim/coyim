package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

func unwrapOrientation(o gtki.Orientation) gtk.Orientation {
	switch o {
	case gtki.HorizontalOrientation:
		return gtk.ORIENTATION_HORIZONTAL
	case gtki.VerticalOrientation:
		return gtk.ORIENTATION_VERTICAL
	}

	panic("This should not happen")
}

func wrapOrientation(o gtk.Orientation) gtki.Orientation {
	switch o {
	case gtk.ORIENTATION_HORIZONTAL:
		return gtki.HorizontalOrientation
	case gtk.ORIENTATION_VERTICAL:
		return gtki.VerticalOrientation
	}

	panic("This should not happen")
}
