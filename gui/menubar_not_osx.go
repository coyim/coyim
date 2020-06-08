// +build !darwin

package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

func getBoolSetting(settings gtki.Settings, name string) bool {
	s, _ := settings.GetProperty(name)
	return s.(bool)
}

func (u *gtkUI) initializeMenus() {
	settings, err := g.gtk.SettingsGetDefault()
	if err != nil {
		panic(err)
	}

	showsAppMenu := getBoolSetting(settings, "gtk-shell-shows-app-menu")
	// showsMenubar := getBoolSetting(settings, "gtk-shell-shows-menubar")

	if showsAppMenu {
		appMenu := u.createSimpleAppMenu()
		u.app.SetAppMenu(appMenu)
	}
}
