package access

import "github.com/coyim/gotk3adapter/gtki"
import "github.com/coyim/gotk3adapter/gdki"

type GTKOSX interface {
	GetApplication() (Application, error)
}

type Application interface {
	Ready()
	SetDockIconPixbuf(gdki.Pixbuf)
	InsertAppMenuItem(gtki.Widget, int)
	SetMenuBar(gtki.MenuShell)
	SetHelpMenu(gtki.MenuItem)
	SetWindowMenu(gtki.MenuItem)
}
