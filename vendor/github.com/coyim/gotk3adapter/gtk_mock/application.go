package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockApplication struct {
	glib_mock.MockApplication
}

func (*MockApplication) GetActiveWindow() gtki.Window {
	return nil
}

func (*MockApplication) AddWindow(gtki.Window)    {}
func (*MockApplication) RemoveWindow(gtki.Window) {}
func (*MockApplication) PrefersAppMenu() bool {
	return false
}

func (*MockApplication) GetAppMenu() glibi.MenuModel {
	return nil
}

func (*MockApplication) SetAppMenu(glibi.MenuModel) {
}

func (*MockApplication) GetMenubar() glibi.MenuModel {
	return nil
}

func (*MockApplication) SetMenubar(glibi.MenuModel) {
}
