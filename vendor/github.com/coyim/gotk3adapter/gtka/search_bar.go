package gtka

import (
	"github.com/coyim/gotk3adapter/gdka"
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type searchBar struct {
	*bin
	internal *gtk.SearchBar
}

func WrapSearchBarSimple(v *gtk.SearchBar) gtki.SearchBar {
	if v == nil {
		return nil
	}
	return &searchBar{WrapBinSimple(&v.Bin).(*bin), v}
}

func WrapSearchBar(v *gtk.SearchBar, e error) (gtki.SearchBar, error) {
	return WrapSearchBarSimple(v), e
}

func UnwrapSearchBar(v gtki.SearchBar) *gtk.SearchBar {
	if v == nil {
		return nil
	}
	return v.(*searchBar).internal
}

func (v *searchBar) ConnectEntry(v1 gtki.Entry) {
	v.internal.ConnectEntry(UnwrapEntry(v1))
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
