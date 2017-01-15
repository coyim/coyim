package i18n

type noLocal struct {
}

var NoLocal = &noLocal{}

func (*noLocal) Local(s string) string {
	return s
}
