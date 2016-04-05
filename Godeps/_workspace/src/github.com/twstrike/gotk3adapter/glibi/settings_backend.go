package glibi

type SettingsBackend interface {
	Object
}

func AssertSettingsBackend(_ SettingsBackend) {}
