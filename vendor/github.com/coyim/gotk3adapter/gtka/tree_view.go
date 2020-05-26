package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type treeView struct {
	*container
	internal *gtk.TreeView
}

func wrapTreeViewSimple(v *gtk.TreeView) *treeView {
	if v == nil {
		return nil
	}
	return &treeView{wrapContainerSimple(&v.Container), v}
}

func wrapTreeView(v *gtk.TreeView, e error) (*treeView, error) {
	return wrapTreeViewSimple(v), e
}

func unwrapTreeView(v gtki.TreeView) *gtk.TreeView {
	if v == nil {
		return nil
	}
	return v.(*treeView).internal
}

func (v *treeView) CollapseRow(v1 gtki.TreePath) bool {
	return v.internal.CollapseRow(unwrapTreePath(v1))
}

func (v *treeView) ExpandAll() {
	v.internal.ExpandAll()
}

func (v *treeView) GetCursor() (gtki.TreePath, gtki.TreeViewColumn) {
	v1, v2 := v.internal.GetCursor()
	return wrapTreePathSimple(v1), wrapTreeViewColumnSimple(v2)
}

func (v *treeView) GetSelection() (gtki.TreeSelection, error) {
	return wrapTreeSelection(v.internal.GetSelection())
}

func (v *treeView) GetPathAtPos(v1 int, v2 int) (gtki.TreePath, gtki.TreeViewColumn, int, int, bool) {
	r1, r2, r3, r4, r5 := v.internal.GetPathAtPos(v1, v2)
	return wrapTreePathSimple(r1), wrapTreeViewColumnSimple(r2), r3, r4, r5
}

func (v *treeView) SetEnableSearch(v1 bool) {
	v.internal.SetEnableSearch(v1)
}

func (v *treeView) GetEnableSearch() bool {
	return v.internal.GetEnableSearch()
}

func (v *treeView) SetSearchColumn(v1 int) {
	v.internal.SetSearchColumn(v1)
}

func (v *treeView) GetSearchColumn() int {
	return v.internal.GetSearchColumn()
}

func (v *treeView) GetSearchEntry() gtki.Entry {
	return wrapEntrySimple(v.internal.GetSearchEntry())
}

func (v *treeView) SetSearchEntry(v1 gtki.Entry) {
	v.internal.SetSearchEntry(unwrapEntry(v1))
}

func (v *treeView) SetSearchEqualSubstringMatch() {
	v.internal.SetSearchEqualSubstringMatch()
}

func (v *treeView) SetModel(m gtki.TreeModel) {
	v.internal.SetModel(unwrapTreeModel(m))
}

func (v *treeView) GetModel() (gtki.TreeModel, error) {
	m, err := v.internal.GetModel()
	if err != nil {
		return nil, err
	}

	return wrapTreeModelSimple(m), nil
}
