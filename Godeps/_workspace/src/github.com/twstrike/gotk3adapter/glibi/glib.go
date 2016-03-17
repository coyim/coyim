package glibi

type Glib interface {
	IdleAdd(interface{}, ...interface{}) (SourceHandle, error)
	InitI18n(string, string)
	Local(string) string
	MainDepth() int
	SignalNew(string) (Signal, error)
} // end of Glib

func AssertGlib(_ Glib) {}
