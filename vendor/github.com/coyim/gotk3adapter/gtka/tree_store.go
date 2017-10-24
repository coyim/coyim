package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type treeStore struct {
	*gliba.Object
	internal *gtk.TreeStore
}

func wrapTreeStoreSimple(v *gtk.TreeStore) *treeStore {
	if v == nil {
		return nil
	}
	return &treeStore{gliba.WrapObjectSimple(v.Object), v}
}

func wrapTreeStore(v *gtk.TreeStore, e error) (*treeStore, error) {
	return wrapTreeStoreSimple(v), e
}

func unwrapTreeStore(v gtki.TreeStore) *gtk.TreeStore {
	if v == nil {
		return nil
	}
	return v.(*treeStore).internal
}

func (v *treeStore) GetIter(v1 gtki.TreePath) (gtki.TreeIter, error) {
	return wrapTreeIter(v.internal.GetIter(unwrapTreePath(v1)))
}

func (v *treeStore) GetIterFirst() (gtki.TreeIter, bool) {
	v1, v2 := v.internal.GetIterFirst()
	return wrapTreeIterSimple(v1), v2
}

func (v *treeStore) GetIterFromString(v1 string) (gtki.TreeIter, error) {
	return wrapTreeIter(v.internal.GetIterFromString(v1))
}

func (v *treeStore) GetPath(v1 gtki.TreeIter) (gtki.TreePath, error) {
	return wrapTreePath(v.internal.GetPath(unwrapTreeIter(v1)))
}

func (v *treeStore) GetValue(v1 gtki.TreeIter, v2 int) (glibi.Value, error) {
	return gliba.WrapValue(v.internal.GetValue(unwrapTreeIter(v1), v2))
}

func (v *treeStore) IterNext(v1 gtki.TreeIter) bool {
	return v.internal.IterNext(unwrapTreeIter(v1))
}

func (v *treeStore) Append(v1 gtki.TreeIter) gtki.TreeIter {
	return wrapTreeIterSimple(v.internal.Append(unwrapTreeIter(v1)))
}

func (v *treeStore) Clear() {
	v.internal.Clear()
}

func (v *treeStore) SetValue(v1 gtki.TreeIter, v2 int, v3 interface{}) error {
	return v.internal.SetValue(unwrapTreeIter(v1), v2, gliba.UnwrapAllGuard(v3))
}
