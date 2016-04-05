package glibi

type SettingsSchema interface {
	Ref() SettingsSchema
	Unref()
	GetID() string
	GetPath() string
	HasKey(string) bool
}

func AssertSettingsSchema(_ SettingsSchema) {}
