package gliba

import "github.com/coyim/gotk3adapter/glibi"
import "github.com/gotk3/gotk3/glib"

type variant struct {
	*glib.Variant
}

func WrapVariant(v *glib.Variant) glibi.Variant {
	if v == nil {
		return nil
	}
	return &variant{v}
}

func UnwrapVariant(v glibi.Variant) *glib.Variant {
	if v == nil {
		return nil
	}
	return v.(*variant).Variant
}

func (v *variant) TypeString() string {
	return v.Variant.TypeString()
}

func (v *variant) IsContainer() bool {
	return v.Variant.IsContainer()
}

func (v *variant) GetBoolean() bool {
	return v.Variant.GetBoolean()
}

func (v *variant) GetString() string {
	return v.Variant.GetString()
}

func (v *variant) GetStrv() []string {
	return v.Variant.GetStrv()
}

func (v *variant) GetInt() (int64, error) {
	return v.Variant.GetInt()
}

func (v *variant) Type() glibi.VariantType {
	return WrapVariantType(v.Variant.Type())
}

func (v *variant) IsType(t glibi.VariantType) bool {
	return v.Variant.IsType(UnwrapVariantType(t))
}

func (v *variant) String() string {
	return v.Variant.String()
}

func (v *variant) AnnotatedString() string {
	return v.Variant.AnnotatedString()
}
