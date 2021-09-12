package gui

import "github.com/coyim/gotk3adapter/gtki"

type mainUserInterface struct {
	mainMenuUserInterface
	mainSearchUserInterface

	app    gtki.Application
	window gtki.ApplicationWindow `gtk-widget:"mainWindow"`

	notificationArea gtki.Box `gtk-widget:"notification-area"`

	tags *tags

	mainBuilder *builder
}

func (m *mainUserInterface) loadUIDefinitions() {
	panicOnDevError(m.mainBuilder.bindObjects(m))
	m.mainMenuUserInterface.loadUIDefinitions(m.mainBuilder)
	m.mainSearchUserInterface.loadUIDefinitions(m.mainBuilder)
}
