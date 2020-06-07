package gliba

import "github.com/coyim/gotk3adapter/glibi"

func init() {
	glibi.AssertGlib(&RealGlib{})
	glibi.AssertApplication(&Application{})
	glibi.AssertObject(&Object{})
	glibi.AssertSettings(&settings{})
	glibi.AssertSettingsBackend(&settingsBackend{})
	glibi.AssertSettingsSchema(&settingsSchema{})
	glibi.AssertSettingsSchemaSource(&settingsSchemaSource{})
	glibi.AssertSignal(&signal{})
	glibi.AssertValue(&value{})
	glibi.AssertMenu(&menu{})
	glibi.AssertMenuItem(&menuItem{})
	glibi.AssertMenuModel(&menuModel{})
	glibi.AssertVariant(&variant{})
	glibi.AssertAction(&action{})
	glibi.AssertSimpleAction(&simpleAction{})
	glibi.AssertPropertyAction(&propertyAction{})
}
