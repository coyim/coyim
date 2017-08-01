package i18n

type noLocal struct {
}

// NoLocal is a i18.Localizer that performs no localization.
var NoLocal = &noLocal{}

func (*noLocal) Local(s string) string {
	return s
}
