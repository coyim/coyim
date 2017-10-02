package gliba

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/coyim/gotk3adapter/glibi"
)

type settingsSchema struct {
	*glib.SettingsSchema
}

func wrapSettingsSchemaSimple(v *glib.SettingsSchema) *settingsSchema {
	if v == nil {
		return nil
	}
	return &settingsSchema{v}
}

func unwrapSettingsSchema(v glibi.SettingsSchema) *glib.SettingsSchema {
	if v == nil {
		return nil
	}
	return v.(*settingsSchema).SettingsSchema
}

func (v *settingsSchema) Ref() glibi.SettingsSchema {
	return wrapSettingsSchemaSimple(v.SettingsSchema.Ref())
}

func (v *settingsSchema) Unref() {
	v.SettingsSchema.Unref()
}

func (v *settingsSchema) GetID() string {
	return v.SettingsSchema.GetID()
}

func (v *settingsSchema) GetPath() string {
	return v.SettingsSchema.GetPath()
}

func (v *settingsSchema) HasKey(v1 string) bool {
	return v.SettingsSchema.HasKey(v1)
}
