package gliba

import "github.com/gotk3/gotk3/glib"

import "github.com/coyim/gotk3adapter/glibi"

type RealGlib struct{}

var Real = &RealGlib{}

func (*RealGlib) IdleAdd(f interface{}, args ...interface{}) (glibi.SourceHandle, error) {
	res, err := glib.IdleAdd(f, args...)
	return glibi.SourceHandle(res), err
}

func (*RealGlib) InitI18n(domain string, dir string) {
	glib.InitI18n(domain, dir)
}

func (*RealGlib) Local(v1 string) string {
	return glib.Local(v1)
}

func (*RealGlib) MainDepth() int {
	return glib.MainDepth()
}

func (*RealGlib) SignalNew(s string) (glibi.Signal, error) {
	return wrapSignal(glib.SignalNew(s))
}

func (*RealGlib) SettingsNew(v1 string) glibi.Settings {
	return WrapSettingsSimple(glib.SettingsNew(v1))
}

func (*RealGlib) SettingsNewWithPath(v1 string, v2 string) glibi.Settings {
	return WrapSettingsSimple(glib.SettingsNewWithPath(v1, v2))
}

func (*RealGlib) SettingsNewWithBackend(v1 string, v2 glibi.SettingsBackend) glibi.Settings {
	return WrapSettingsSimple(glib.SettingsNewWithBackend(v1, UnwrapSettingsBackend(v2)))
}

func (*RealGlib) SettingsNewWithBackendAndPath(v1 string, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	return WrapSettingsSimple(glib.SettingsNewWithBackendAndPath(v1, UnwrapSettingsBackend(v2), v3))
}

func (*RealGlib) SettingsNewFull(v1 glibi.SettingsSchema, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	return WrapSettingsSimple(glib.SettingsNewFull(UnwrapSettingsSchema(v1), UnwrapSettingsBackend(v2), v3))
}

func (*RealGlib) SettingsSync() {
	glib.SettingsSync()
}

func (*RealGlib) SettingsBackendGetDefault() glibi.SettingsBackend {
	return WrapSettingsBackendSimple(glib.SettingsBackendGetDefault())
}

func (*RealGlib) KeyfileSettingsBackendNew(v1 string, v2 string, v3 string) glibi.SettingsBackend {
	return WrapSettingsBackendSimple(glib.KeyfileSettingsBackendNew(v1, v2, v3))
}

func (*RealGlib) MemorySettingsBackendNew() glibi.SettingsBackend {
	return WrapSettingsBackendSimple(glib.MemorySettingsBackendNew())
}

func (*RealGlib) NullSettingsBackendNew() glibi.SettingsBackend {
	return WrapSettingsBackendSimple(glib.NullSettingsBackendNew())
}

func (*RealGlib) SettingsSchemaSourceGetDefault() glibi.SettingsSchemaSource {
	return WrapSettingsSchemaSourceSimple(glib.SettingsSchemaSourceGetDefault())
}

func (*RealGlib) SettingsSchemaSourceNewFromDirectory(v1 string, v2 glibi.SettingsSchemaSource, v3 bool) glibi.SettingsSchemaSource {
	return WrapSettingsSchemaSourceSimple(glib.SettingsSchemaSourceNewFromDirectory(v1, UnwrapSettingsSchemaSource(v2), v3))
}

func (*RealGlib) MenuNew() glibi.Menu {
	return WrapMenuSimple(glib.MenuNew())
}

func (*RealGlib) MenuItemNew(label, detailed_action string) glibi.MenuItem {
	return WrapMenuItemSimple(glib.MenuItemNew(label, detailed_action))
}

func (*RealGlib) MenuItemNewSection(label string, section glibi.MenuModel) glibi.MenuItem {
	return WrapMenuItemSimple(glib.MenuItemNewSection(label, UnwrapMenuModel(section)))
}

func (*RealGlib) MenuItemNewSubmenu(label string, submenu glibi.MenuModel) glibi.MenuItem {
	return WrapMenuItemSimple(glib.MenuItemNewSubmenu(label, UnwrapMenuModel(submenu)))
}

func (*RealGlib) MenuItemNewFromModel(model glibi.MenuModel, index int) glibi.MenuItem {
	return WrapMenuItemSimple(glib.MenuItemNewFromModel(UnwrapMenuModel(model), index))
}
