package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type textTagTable struct {
	*gliba.Object
	internal *gtk.TextTagTable
}

func wrapTextTagTableSimple(v *gtk.TextTagTable) *textTagTable {
	if v == nil {
		return nil
	}
	return &textTagTable{gliba.WrapObjectSimple(v.Object), v}
}

func wrapTextTagTable(v *gtk.TextTagTable, e error) (*textTagTable, error) {
	return wrapTextTagTableSimple(v), e
}

func unwrapTextTagTable(v gtki.TextTagTable) *gtk.TextTagTable {
	if v == nil {
		return nil
	}
	return v.(*textTagTable).internal
}

func (v *textTagTable) Add(v1 gtki.TextTag) {
	v.internal.Add(unwrapTextTag(v1))
}
