package pangoa

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/pango"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/pangoi"
)

func init() {
	pangoi.PANGO_SCALE = pango.PANGO_SCALE
}
