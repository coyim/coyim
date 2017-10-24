package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
)

type searchBar struct {
	*bin
	internal *gtk.SearchBar
}

func wrapSearchBarSimple(v *gtk.SearchBar) *searchBar {
	if v == nil {
		return nil
	}
	return &searchBar{wrapBinSimple(&v.Bin), v}
}

func wrapSearchBar(v *gtk.SearchBar, e error) (*searchBar, error) {
	return wrapSearchBarSimple(v), e
}

func unwrapSearchBar(v gtki.SearchBar) *gtk.SearchBar {
	if v == nil {
		return nil
	}
	return v.(*searchBar).internal
}

func (v *searchBar) ConnectEntry(v1 gtki.Entry) {
	v.internal.ConnectEntry(unwrapEntry(v1))
}

func (v *searchBar) GetSearchMode() bool {
	return v.internal.GetSearchMode()
}

func (v *searchBar) SetSearchMode(v1 bool) {
	v.internal.SetSearchMode(v1)
}

func (v *searchBar) GetShowCloseButton() bool {
	return v.internal.GetShowCloseButton()
}

func (v *searchBar) SetShowCloseButton(v1 bool) {
	v.internal.SetShowCloseButton(v1)
}

func (v *searchBar) HandleEvent(v1 gdki.Event) {
	v.internal.HandleEvent(gdka.UnwrapEvent(v1))
}
