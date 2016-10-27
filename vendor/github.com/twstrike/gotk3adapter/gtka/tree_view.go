package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/gotk3adapter/gtki"
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

func (v *treeView) GetPathAtPos(v1 int, v2 int, v3 gtki.TreePath, v4 gtki.TreeViewColumn, v5 *int, v6 *int) bool {
	return v.internal.GetPathAtPos(v1, v2, unwrapTreePath(v3), unwrapTreeViewColumn(v4), v5, v6)
}
