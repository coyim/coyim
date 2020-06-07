package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type searchEntry struct {
	*entry
	internal *gtk.SearchEntry
}

func WrapSearchEntrySimple(v *gtk.SearchEntry) gtki.SearchEntry {
	if v == nil {
		return nil
	}
	return &searchEntry{WrapEntrySimple(&v.Entry).(*entry), v}
}

func WrapSearchEntry(v *gtk.SearchEntry, e error) (gtki.SearchEntry, error) {
	return WrapSearchEntrySimple(v), e
}

func UnwrapSearchEntry(v gtki.SearchEntry) *gtk.SearchEntry {
	if v == nil {
		return nil
	}
	return v.(*searchEntry).internal
}
