package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type treeSelection struct {
	*gliba.Object
	internal *gtk.TreeSelection
}

func wrapTreeSelectionSimple(v *gtk.TreeSelection) *treeSelection {
	if v == nil {
		return nil
	}
	return &treeSelection{gliba.WrapObjectSimple(v.Object), v}
}

func wrapTreeSelection(v *gtk.TreeSelection, e error) (*treeSelection, error) {
	return wrapTreeSelectionSimple(v), e
}

func unwrapTreeSelection(v gtki.TreeSelection) *gtk.TreeSelection {
	if v == nil {
		return nil
	}
	return v.(*treeSelection).internal
}

func (v *treeSelection) SelectIter(v1 gtki.TreeIter) {
	v.internal.SelectIter(unwrapTreeIter(v1))
}

func (v *treeSelection) UnselectPath(v1 gtki.TreePath) {
	v.internal.UnselectPath(unwrapTreePath(v1))
}

func (v *treeSelection) GetSelected() (gtki.TreeModel, gtki.TreeIter, bool) {
	v1, v2, v3 := v.internal.GetSelected()
	return wrapTreeModelSimple(v1), wrapTreeIterSimple(v2), v3
}
