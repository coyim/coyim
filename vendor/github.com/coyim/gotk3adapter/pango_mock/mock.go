package pango_mock

import "github.com/coyim/gotk3adapter/pangoi"

type Mock struct{}

func (*Mock) AsFontDescription(v interface{}) pangoi.FontDescription {
	return nil
}
