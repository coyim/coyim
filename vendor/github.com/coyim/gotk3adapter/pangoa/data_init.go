package pangoa

import (
	"github.com/coyim/gotk3adapter/pangoi"
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
}
