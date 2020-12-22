package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type treeStore struct {
	*gliba.Object
	internal *gtk.TreeStore
}

func WrapTreeStoreSimple(v *gtk.TreeStore) gtki.TreeStore {
	if v == nil {
		return nil
	}
	return &treeStore{gliba.WrapObjectSimple(v.Object), v}
}

func WrapTreeStore(v *gtk.TreeStore, e error) (gtki.TreeStore, error) {
	return WrapTreeStoreSimple(v), e
}

func UnwrapTreeStore(v gtki.TreeStore) *gtk.TreeStore {
	if v == nil {
		return nil
	}
	return v.(*treeStore).internal
}

func (v *treeStore) GetIter(v1 gtki.TreePath) (gtki.TreeIter, error) {
	return WrapTreeIter(v.internal.GetIter(UnwrapTreePath(v1)))
}

func (v *treeStore) GetIterFirst() (gtki.TreeIter, bool) {
	v1, v2 := v.internal.GetIterFirst()
	return WrapTreeIterSimple(v1), v2
}

func (v *treeStore) GetIterFromString(v1 string) (gtki.TreeIter, error) {
	return WrapTreeIter(v.internal.GetIterFromString(v1))
}

func (v *treeStore) GetPath(v1 gtki.TreeIter) (gtki.TreePath, error) {
	return WrapTreePath(v.internal.GetPath(UnwrapTreeIter(v1)))
}

func (v *treeStore) GetValue(v1 gtki.TreeIter, v2 int) (glibi.Value, error) {
	return gliba.WrapValue(v.internal.GetValue(UnwrapTreeIter(v1), v2))
}

func (v *treeStore) IterNext(v1 gtki.TreeIter) bool {
	return v.internal.IterNext(UnwrapTreeIter(v1))
}

func (v *treeStore) Append(v1 gtki.TreeIter) gtki.TreeIter {
	return WrapTreeIterSimple(v.internal.Append(UnwrapTreeIter(v1)))
}

func (v *treeStore) Clear() {
	v.internal.Clear()
}

func (v *treeStore) SetValue(v1 gtki.TreeIter, v2 int, v3 interface{}) error {
	return v.internal.SetValue(UnwrapTreeIter(v1), v2, gliba.UnwrapAllGuard(v3))
}

func (v *treeStore) Remove(v1 gtki.TreeIter) bool {
	return v.internal.Remove(UnwrapTreeIter(v1))
}
