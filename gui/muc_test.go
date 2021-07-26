package gui

import (
	"fmt"

	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"
)

type mucMockGlib struct {
	glib_mock.Mock
}

func (*mucMockGlib) Local(vx string) string {
	return "[localized] " + vx
}

func (*mucMockGlib) Localf(vx string, args ...interface{}) string {
	return fmt.Sprintf("[localized] "+vx, args...)
}

func initMUCi18n() {
	i18n.InitLocalization(&mucMockGlib{})
}
