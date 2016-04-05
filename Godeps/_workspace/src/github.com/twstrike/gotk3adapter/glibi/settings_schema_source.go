package glibi

type SettingsSchemaSource interface {
	Ref() SettingsSchemaSource
	Unref()
	Lookup(string, bool) SettingsSchema
}

func AssertSettingsSchemaSource(_ SettingsSchemaSource) {}
