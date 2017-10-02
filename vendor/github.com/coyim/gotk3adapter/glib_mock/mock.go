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
