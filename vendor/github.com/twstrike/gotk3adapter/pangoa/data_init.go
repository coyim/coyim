package pangoa

import (
	"github.com/gotk3/gotk3/pango"
	"github.com/twstrike/gotk3adapter/pangoi"
)

func init() {
	pangoi.PANGO_SCALE = pango.PANGO_SCALE

	pangoi.UNDERLINE_NONE = int(pango.UNDERLINE_NONE)
	pangoi.UNDERLINE_SINGLE = int(pango.UNDERLINE_SINGLE)
	pangoi.UNDERLINE_DOUBLE = int(pango.UNDERLINE_DOUBLE)
	pangoi.UNDERLINE_LOW = int(pango.UNDERLINE_LOW)
	pangoi.UNDERLINE_ERROR = int(pango.UNDERLINE_ERROR)
}
