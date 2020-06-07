package glib_mock

import "github.com/coyim/gotk3adapter/glibi"

type Mock struct{}

func (*Mock) IdleAdd(f interface{}, args ...interface{}) (glibi.SourceHandle, error) {
	return glibi.SourceHandle(0), nil
}

func (*Mock) InitI18n(domain string, dir string) {
}

func (*Mock) Local(vx string) string {
	return vx
}

func (*Mock) MainDepth() int {
	return 0
}

func (*Mock) SignalNew(s string) (glibi.Signal, error) {
	return &MockSignal{}, nil
}

func (*Mock) SettingsNew(string) glibi.Settings {
	return nil
}

func (*Mock) SettingsNewWithPath(string, string) glibi.Settings {
	return nil
}

func (*Mock) SettingsNewWithBackend(string, glibi.SettingsBackend) glibi.Settings {
	return nil
}

func (*Mock) SettingsNewWithBackendAndPath(string, glibi.SettingsBackend, string) glibi.Settings {
	return nil
}

func (*Mock) SettingsNewFull(glibi.SettingsSchema, glibi.SettingsBackend, string) glibi.Settings {
	return nil
}

func (*Mock) SettingsSync() {
}

func (*Mock) SettingsBackendGetDefault() glibi.SettingsBackend {
	return nil
}

func (*Mock) KeyfileSettingsBackendNew(string, string, string) glibi.SettingsBackend {
	return nil
}

func (*Mock) MemorySettingsBackendNew() glibi.SettingsBackend {
	return nil
}

func (*Mock) NullSettingsBackendNew() glibi.SettingsBackend {
	return nil
}

func (*Mock) SettingsSchemaSourceGetDefault() glibi.SettingsSchemaSource {
	return nil
}

func (*Mock) SettingsSchemaSourceNewFromDirectory(string, glibi.SettingsSchemaSource, bool) glibi.SettingsSchemaSource {
	return nil
}

func (*Mock) MenuNew() glibi.Menu {
	return nil
}

func (*Mock) MenuItemNew(label, detailed_action string) glibi.MenuItem {
	return nil
}

func (*Mock) MenuItemNewSection(label string, section glibi.MenuModel) glibi.MenuItem {
	return nil
}

func (*Mock) MenuItemNewSubmenu(label string, submenu glibi.MenuModel) glibi.MenuItem {
	return nil
}

func (*Mock) MenuItemNewFromModel(model glibi.MenuModel, index int) glibi.MenuItem {
	return nil
}

func (*Mock) ActionNameIsValid(actionName string) bool {
	return false
}

func (*Mock) SimpleActionNew(name string, parameterType glibi.VariantType) glibi.SimpleAction {
	return nil
}

func (*Mock) SimpleActionNewStateful(name string, parameterType glibi.VariantType, state glibi.Variant) glibi.SimpleAction {
	return nil
}

func (*Mock) PropertyActionNew(name string, object glibi.Object, propertyName string) glibi.PropertyAction {
	return nil
}
