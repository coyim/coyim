package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type textTagTable struct {
	*gliba.Object
	internal *gtk.TextTagTable
}

func WrapTextTagTableSimple(v *gtk.TextTagTable) gtki.TextTagTable {
	if v == nil {
		return nil
	}
	return &textTagTable{gliba.WrapObjectSimple(v.Object), v}
}

func WrapTextTagTable(v *gtk.TextTagTable, e error) (gtki.TextTagTable, error) {
	return WrapTextTagTableSimple(v), e
}

func UnwrapTextTagTable(v gtki.TextTagTable) *gtk.TextTagTable {
	if v == nil {
		return nil
	}
	return v.(*textTagTable).internal
}

func (v *textTagTable) Add(v1 gtki.TextTag) {
	v.internal.Add(UnwrapTextTag(v1))
}
