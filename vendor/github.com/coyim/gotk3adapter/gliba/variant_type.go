package gliba

import "github.com/coyim/gotk3adapter/glibi"
import "github.com/gotk3/gotk3/glib"

type VariantType struct {
	*glib.VariantType
}

func WrapVariantType(v *glib.VariantType) glibi.VariantType {
	if v == nil {
		return nil
	}
	return &VariantType{v}
}

func UnwrapVariantType(v glibi.VariantType) *glib.VariantType {
	if v == nil {
		return nil
	}
	return v.(*VariantType).VariantType
}
