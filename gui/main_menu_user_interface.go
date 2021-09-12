package gui

import "github.com/coyim/gotk3adapter/gtki"

type mainMenuUserInterface struct {
	contactsMenuItem  gtki.MenuItem `gtk-widget:"ContactsMenu"`
	accountsMenuItem  gtki.MenuItem `gtk-widget:"AccountsMenu"`
	chatRoomsMenuItem gtki.MenuItem `gtk-widget:"ChatRoomsMenu"`
	viewMenuItem      gtki.MenuItem `gtk-widget:"ViewMenu"`
	optionsMenuItem   gtki.MenuItem `gtk-widget:"OptionsMenu"`

	viewMenu    *viewMenu
	optionsMenu *optionsMenu

	toggleConnectAllAutomaticallyRequest chan bool
	setShowAdvancedSettingsRequest       chan bool
}

func (m *mainMenuUserInterface) loadUIDefinitions(b *builder) {
	panicOnDevError(b.bindObjects(m))
}
