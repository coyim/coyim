package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type settingsSchemaSource struct {
	*glib.SettingsSchemaSource
}

func WrapSettingsSchemaSourceSimple(v *glib.SettingsSchemaSource) glibi.SettingsSchemaSource {
	if v == nil {
		return nil
	}
	return &settingsSchemaSource{v}
}

func UnwrapSettingsSchemaSource(v glibi.SettingsSchemaSource) *glib.SettingsSchemaSource {
	if v == nil {
		return nil
	}
	return v.(*settingsSchemaSource).SettingsSchemaSource
}

func (v *settingsSchemaSource) Ref() glibi.SettingsSchemaSource {
	return WrapSettingsSchemaSourceSimple(v.SettingsSchemaSource.Ref())
}

func (v *settingsSchemaSource) Unref() {
	v.SettingsSchemaSource.Unref()
}

func (v *settingsSchemaSource) Lookup(v1 string, v2 bool) glibi.SettingsSchema {
	return WrapSettingsSchemaSimple(v.SettingsSchemaSource.Lookup(v1, v2))
}
