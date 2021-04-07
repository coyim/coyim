package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type treeSelection struct {
	*gliba.Object
	internal *gtk.TreeSelection
}

func WrapTreeSelectionSimple(v *gtk.TreeSelection) gtki.TreeSelection {
	if v == nil {
		return nil
	}
	return &treeSelection{gliba.WrapObjectSimple(v.Object), v}
}

func WrapTreeSelection(v *gtk.TreeSelection, e error) (gtki.TreeSelection, error) {
	return WrapTreeSelectionSimple(v), e
}

func UnwrapTreeSelection(v gtki.TreeSelection) *gtk.TreeSelection {
	if v == nil {
		return nil
	}
	return v.(*treeSelection).internal
}

func (v *treeSelection) SelectIter(v1 gtki.TreeIter) {
	v.internal.SelectIter(UnwrapTreeIter(v1))
}

func (v *treeSelection) UnselectPath(v1 gtki.TreePath) {
	v.internal.UnselectPath(UnwrapTreePath(v1))
}

func (v *treeSelection) GetSelected() (gtki.TreeModel, gtki.TreeIter, bool) {
	v1, v2, v3 := v.internal.GetSelected()
	return WrapTreeModelSimple(v1), WrapTreeIterSimple(v2), v3
}

func (v *treeSelection) GetSelectedRows(m gtki.TreeModel) []gtki.TreePath {
	ll := v.internal.GetSelectedRows(UnwrapTreeModel(m))

	result := []gtki.TreePath{}
	for cc := ll; cc != nil; cc = cc.Next() {
		result = append(result, Wrap(cc.Data()).(gtki.TreePath))
	}

	return result
}
