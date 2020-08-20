package i18n

import "fmt"

// Localizer represents something that can localize a string
type Localizer interface {
	Local(string) string
}

type nullLocalizer struct{}

func (*nullLocalizer) Local(v string) string {
	return fmt.Sprintf("[NULL LOCALIZER] - %s", v)
}

var g Localizer = &nullLocalizer{}

// InitLocalization should be called before using localization - it sets the variable used to access the localization interface
func InitLocalization(gx Localizer) {
	g = gx
}

// Local returns the given string in the local language
func Local(v string) string {
	return g.Local(v)
}

// Localf returns the given string in the local language. It supports Printf formatting.
func Localf(f string, p ...interface{}) string {
	return fmt.Sprintf(Local(f), p...)
}
