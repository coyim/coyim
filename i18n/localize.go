package i18n

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
)

// Local returns the given string in the local language
func Local(v string) string {
	return glib.Local(v)
}

// Localf returns the given string in the local language. It supports Printf formatting.
func Localf(f string, p ...interface{}) string {
	return fmt.Sprintf(Local(f), p...)
}
