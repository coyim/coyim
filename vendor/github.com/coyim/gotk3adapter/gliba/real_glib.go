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
	return wrapSettingsSimple(glib.SettingsNew(v1))
}

func (*RealGlib) SettingsNewWithPath(v1 string, v2 string) glibi.Settings {
	return wrapSettingsSimple(glib.SettingsNewWithPath(v1, v2))
}

func (*RealGlib) SettingsNewWithBackend(v1 string, v2 glibi.SettingsBackend) glibi.Settings {
	return wrapSettingsSimple(glib.SettingsNewWithBackend(v1, unwrapSettingsBackend(v2)))
}

func (*RealGlib) SettingsNewWithBackendAndPath(v1 string, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	return wrapSettingsSimple(glib.SettingsNewWithBackendAndPath(v1, unwrapSettingsBackend(v2), v3))
}

func (*RealGlib) SettingsNewFull(v1 glibi.SettingsSchema, v2 glibi.SettingsBackend, v3 string) glibi.Settings {
	return wrapSettingsSimple(glib.SettingsNewFull(unwrapSettingsSchema(v1), unwrapSettingsBackend(v2), v3))
}

func (*RealGlib) SettingsSync() {
	glib.SettingsSync()
}

func (*RealGlib) SettingsBackendGetDefault() glibi.SettingsBackend {
	return wrapSettingsBackendSimple(glib.SettingsBackendGetDefault())
}

func (*RealGlib) KeyfileSettingsBackendNew(v1 string, v2 string, v3 string) glibi.SettingsBackend {
	return wrapSettingsBackendSimple(glib.KeyfileSettingsBackendNew(v1, v2, v3))
}

func (*RealGlib) MemorySettingsBackendNew() glibi.SettingsBackend {
	return wrapSettingsBackendSimple(glib.MemorySettingsBackendNew())
}

func (*RealGlib) NullSettingsBackendNew() glibi.SettingsBackend {
	return wrapSettingsBackendSimple(glib.NullSettingsBackendNew())
}

func (*RealGlib) SettingsSchemaSourceGetDefault() glibi.SettingsSchemaSource {
	return wrapSettingsSchemaSourceSimple(glib.SettingsSchemaSourceGetDefault())
}

func (*RealGlib) SettingsSchemaSourceNewFromDirectory(v1 string, v2 glibi.SettingsSchemaSource, v3 bool) glibi.SettingsSchemaSource {
	return wrapSettingsSchemaSourceSimple(glib.SettingsSchemaSourceNewFromDirectory(v1, unwrapSettingsSchemaSource(v2), v3))
}
