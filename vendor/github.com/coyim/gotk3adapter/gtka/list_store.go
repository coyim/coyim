package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type listStore struct {
	*gliba.Object
	internal *gtk.ListStore
}

func wrapListStoreSimple(v *gtk.ListStore) *listStore {
	if v == nil {
		return nil
	}
	return &listStore{gliba.WrapObjectSimple(v.Object), v}
}

func wrapListStore(v *gtk.ListStore, e error) (*listStore, error) {
	return wrapListStoreSimple(v), e
}

func unwrapListStore(v gtki.ListStore) *gtk.ListStore {
	if v == nil {
		return nil
	}
	return v.(*listStore).internal
}

func (v *listStore) Clear() {
	v.internal.Clear()
}

func (v *listStore) Append() gtki.TreeIter {
	return wrapTreeIterSimple(v.internal.Append())
}

func (v *listStore) Remove(v2 gtki.TreeIter) bool {
	return v.internal.Remove(unwrapTreeIter(v2))
}

func (v *listStore) Set2(v2 gtki.TreeIter, v3 []int, v4 []interface{}) error {
	return v.internal.Set(unwrapTreeIter(v2), v3, gliba.FixupArray(v4))
}

func (v *listStore) SetValue(v1 gtki.TreeIter, v2 int, v3 interface{}) error {
	return v.internal.SetValue(unwrapTreeIter(v1), v2, gliba.UnwrapAllGuard(v3))
}

func (v *listStore) GetIter(v1 gtki.TreePath) (gtki.TreeIter, error) {
	return wrapTreeIter(v.internal.GetIter(unwrapTreePath(v1)))
}

func (v *listStore) GetIterFirst() (gtki.TreeIter, bool) {
	v1, v2 := v.internal.GetIterFirst()
	return wrapTreeIterSimple(v1), v2
}

func (v *listStore) GetIterFromString(v1 string) (gtki.TreeIter, error) {
	return wrapTreeIter(v.internal.GetIterFromString(v1))
}

func (v *listStore) GetPath(v1 gtki.TreeIter) (gtki.TreePath, error) {
	return wrapTreePath(v.internal.GetPath(unwrapTreeIter(v1)))
}

func (v *listStore) GetValue(v1 gtki.TreeIter, v2 int) (glibi.Value, error) {
	return gliba.WrapValue(v.internal.GetValue(unwrapTreeIter(v1), v2))
}

func (v *listStore) IterNext(v1 gtki.TreeIter) bool {
	return v.internal.IterNext(unwrapTreeIter(v1))
}
