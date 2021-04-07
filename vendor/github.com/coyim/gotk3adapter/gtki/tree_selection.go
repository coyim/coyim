package gtki

import "github.com/coyim/gotk3adapter/glibi"

type TreeSelection interface {
	glibi.Object

	GetSelected() (TreeModel, TreeIter, bool)
	GetSelectedRows(TreeModel) []TreePath
	SelectIter(TreeIter)
	UnselectPath(TreePath)
}

func AssertTreeSelection(_ TreeSelection) {}
