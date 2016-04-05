package glib_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type MockSettingsSchema struct {
}

func (*MockSettingsSchema) Ref() glibi.SettingsSchema {
	return nil
}

func (*MockSettingsSchema) Unref() {}

func (*MockSettingsSchema) GetID() string {
	return ""
}

func (*MockSettingsSchema) GetPath() string {
	return ""
}

func (*MockSettingsSchema) HasKey(string) bool {
	return false
}
