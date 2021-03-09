package pango_mock

import "github.com/coyim/gotk3adapter/pangoi"

func init() {
	pangoi.AssertPango(&Mock{})
	pangoi.AssertFontDescription(&MockFontDescription{})
}
