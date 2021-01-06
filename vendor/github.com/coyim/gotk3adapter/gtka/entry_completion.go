package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type entryCompletion struct {
	*gliba.Object
	internal *gtk.EntryCompletion
}

func WrapEntryCompletionSimple(v *gtk.EntryCompletion) gtki.EntryCompletion {
	if v == nil {
		return nil
	}
	return &entryCompletion{gliba.WrapObjectSimple(v.Object), v}
}

func WrapEntryCompletion(v *gtk.EntryCompletion, e error) (gtki.EntryCompletion, error) {
	return WrapEntryCompletionSimple(v), e
}

func UnwrapEntryCompletion(v gtki.EntryCompletion) *gtk.EntryCompletion {
	if v == nil {
		return nil
	}
	return v.(*entryCompletion).internal
}

func (v *entryCompletion) SetModel(v1 gtki.TreeModel) {
	v.internal.SetModel(UnwrapTreeModel(v1))
}

func (v *entryCompletion) SetTextColumn(v1 int) {
	v.internal.SetTextColumn(v1)
}

func (v *entryCompletion) SetMinimumKeyLength(v1 int) {
	v.internal.SetMinimumKeyLength(v1)
}
