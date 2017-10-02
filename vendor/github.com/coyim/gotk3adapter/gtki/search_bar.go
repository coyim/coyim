package gtki

import "github.com/coyim/gotk3adapter/gdki"

type SearchBar interface {
	Bin

	ConnectEntry(Entry)
	GetSearchMode() bool
	SetSearchMode(bool)
	GetShowCloseButton() bool
	SetShowCloseButton(bool)
	HandleEvent(gdki.Event)
}

func AssertSearchBar(_ SearchBar) {}
