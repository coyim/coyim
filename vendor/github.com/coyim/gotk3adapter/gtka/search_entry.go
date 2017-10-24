package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type searchEntry struct {
	*entry
	internal *gtk.SearchEntry
}

func wrapSearchEntrySimple(v *gtk.SearchEntry) *searchEntry {
	if v == nil {
		return nil
	}
	return &searchEntry{wrapEntrySimple(&v.Entry), v}
}

func wrapSearchEntry(v *gtk.SearchEntry, e error) (*searchEntry, error) {
	return wrapSearchEntrySimple(v), e
}

func unwrapSearchEntry(v gtki.SearchEntry) *gtk.SearchEntry {
	if v == nil {
		return nil
	}
	return v.(*searchEntry).internal
}
