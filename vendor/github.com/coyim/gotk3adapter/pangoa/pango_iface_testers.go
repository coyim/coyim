package pangoa

import "github.com/coyim/gotk3adapter/pangoi"

func init() {
	pangoi.AssertPango(&RealPango{})
	pangoi.AssertFontDescription(&fontDescription{})
}
