package glibi

type Glib interface {
	IdleAdd(interface{}, ...interface{}) (SourceHandle, error)
	InitI18n(string, string)
	Local(string) string
	MainDepth() int

	SettingsNew(string) Settings
	SettingsNewWithPath(string, string) Settings
	SettingsNewWithBackend(string, SettingsBackend) Settings
	SettingsNewWithBackendAndPath(string, SettingsBackend, string) Settings
	SettingsNewFull(SettingsSchema, SettingsBackend, string) Settings
	SettingsSync()

	SettingsBackendGetDefault() SettingsBackend
	KeyfileSettingsBackendNew(string, string, string) SettingsBackend
	MemorySettingsBackendNew() SettingsBackend
	NullSettingsBackendNew() SettingsBackend

	SettingsSchemaSourceGetDefault() SettingsSchemaSource
	SettingsSchemaSourceNewFromDirectory(string, SettingsSchemaSource, bool) SettingsSchemaSource

	SignalNew(string) (Signal, error)

	MenuNew() Menu
	MenuItemNew(label, detailed_action string) MenuItem
	MenuItemNewSection(label string, section MenuModel) MenuItem
	MenuItemNewSubmenu(label string, submenu MenuModel) MenuItem
	MenuItemNewFromModel(model MenuModel, index int) MenuItem
} // end of Glib

func AssertGlib(_ Glib) {}
