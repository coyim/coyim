package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type listStore struct {
	*gliba.Object
	internal *gtk.ListStore
}

func WrapListStoreSimple(v *gtk.ListStore) gtki.ListStore {
	if v == nil {
		return nil
	}
	return &listStore{gliba.WrapObjectSimple(v.Object), v}
}

func WrapListStore(v *gtk.ListStore, e error) (gtki.ListStore, error) {
	return WrapListStoreSimple(v), e
}

func UnwrapListStore(v gtki.ListStore) *gtk.ListStore {
	if v == nil {
		return nil
	}
	return v.(*listStore).internal
}

func (v *listStore) Clear() {
	v.internal.Clear()
}

func (v *listStore) Append() gtki.TreeIter {
	return WrapTreeIterSimple(v.internal.Append())
}

func (v *listStore) Remove(v2 gtki.TreeIter) bool {
	return v.internal.Remove(UnwrapTreeIter(v2))
}

func (v *listStore) Set2(v2 gtki.TreeIter, v3 []int, v4 []interface{}) error {
	return v.internal.Set(UnwrapTreeIter(v2), v3, gliba.FixupArray(v4))
}

func (v *listStore) SetValue(v1 gtki.TreeIter, v2 int, v3 interface{}) error {
	return v.internal.SetValue(UnwrapTreeIter(v1), v2, gliba.UnwrapAllGuard(v3))
}

func (v *listStore) GetIter(v1 gtki.TreePath) (gtki.TreeIter, error) {
	return WrapTreeIter(v.internal.GetIter(UnwrapTreePath(v1)))
}

func (v *listStore) GetIterFirst() (gtki.TreeIter, bool) {
	v1, v2 := v.internal.GetIterFirst()
	return WrapTreeIterSimple(v1), v2
}

func (v *listStore) GetIterFromString(v1 string) (gtki.TreeIter, error) {
	return WrapTreeIter(v.internal.GetIterFromString(v1))
}

func (v *listStore) GetPath(v1 gtki.TreeIter) (gtki.TreePath, error) {
	return WrapTreePath(v.internal.GetPath(UnwrapTreeIter(v1)))
}

func (v *listStore) GetValue(v1 gtki.TreeIter, v2 int) (glibi.Value, error) {
	return gliba.WrapValue(v.internal.GetValue(UnwrapTreeIter(v1), v2))
}

func (v *listStore) IterNext(v1 gtki.TreeIter) bool {
	return v.internal.IterNext(UnwrapTreeIter(v1))
}

func (v *listStore) GetColumnType(index int) glibi.Type {
	return glibi.Type(v.internal.GetColumnType(index))
}
