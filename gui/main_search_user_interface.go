package gui

import "github.com/coyim/gotk3adapter/gtki"

type mainSearchUserInterface struct {
	searchBox   gtki.Box       `gtk-widget:"search-box"`
	search      gtki.SearchBar `gtk-widget:"search-area"`
	searchEntry gtki.Entry     `gtk-widget:"search-entry"`
}

func (m *mainSearchUserInterface) loadUIDefinitions(b *builder) {
	panicOnDevError(b.bindObjects(m))
}
