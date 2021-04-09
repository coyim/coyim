package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

func init() {
	pangoi.PANGO_SCALE = pango.PANGO_SCALE

	pangoi.UNDERLINE_NONE = int(pango.UNDERLINE_NONE)
	pangoi.UNDERLINE_SINGLE = int(pango.UNDERLINE_SINGLE)
	pangoi.UNDERLINE_DOUBLE = int(pango.UNDERLINE_DOUBLE)
	pangoi.UNDERLINE_LOW = int(pango.UNDERLINE_LOW)
	pangoi.UNDERLINE_ERROR = int(pango.UNDERLINE_ERROR)

	pangoi.STYLE_NORMAL = int(pango.STYLE_NORMAL)
	pangoi.STYLE_OBLIQUE = int(pango.STYLE_OBLIQUE)
	pangoi.STYLE_ITALIC = int(pango.STYLE_ITALIC)

	pangoi.JUSTIFY_LEFT = int(gtk.JUSTIFY_LEFT)
	pangoi.JUSTIFY_RIGHT = int(gtk.JUSTIFY_RIGHT)
	pangoi.JUSTIFY_CENTER = int(gtk.JUSTIFY_CENTER)
	pangoi.JUSTIFY_FILL = int(gtk.JUSTIFY_FILL)

	pangoi.WEIGHT_THIN = int(pango.WEIGHT_THIN)
	pangoi.WEIGHT_ULTRALIGHT = int(pango.WEIGHT_ULTRALIGHT)
	pangoi.WEIGHT_LIGHT = int(pango.WEIGHT_LIGHT)
	pangoi.WEIGHT_SEMILIGHT = int(pango.WEIGHT_SEMILIGHT)
	pangoi.WEIGHT_BOOK = int(pango.WEIGHT_BOOK)
	pangoi.WEIGHT_NORMAL = int(pango.WEIGHT_NORMAL)
	pangoi.WEIGHT_MEDIUM = int(pango.WEIGHT_MEDIUM)
	pangoi.WEIGHT_SEMIBOLD = int(pango.WEIGHT_SEMIBOLD)
	pangoi.WEIGHT_BOLD = int(pango.WEIGHT_BOLD)
	pangoi.WEIGHT_ULTRABOLD = int(pango.WEIGHT_ULTRABOLD)
	pangoi.WEIGHT_HEAVY = int(pango.WEIGHT_HEAVY)
	pangoi.WEIGHT_ULTRAHEAVY = int(pango.WEIGHT_ULTRAHEAVY)
}
