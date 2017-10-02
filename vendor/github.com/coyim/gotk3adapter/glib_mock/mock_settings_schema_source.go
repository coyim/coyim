package glib_mock

import "github.com/coyim/gotk3adapter/glibi"

type MockSettingsSchemaSource struct {
}

func (*MockSettingsSchemaSource) Ref() glibi.SettingsSchemaSource {
	return nil
}

func (*MockSettingsSchemaSource) Unref() {}

func (*MockSettingsSchemaSource) Lookup(string, bool) glibi.SettingsSchema {
	return nil
}
