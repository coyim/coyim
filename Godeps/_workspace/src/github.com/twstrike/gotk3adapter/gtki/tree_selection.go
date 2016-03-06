package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type TreeSelection interface {
	glibi.Object

	GetSelected() (TreeModel, TreeIter, bool)
	SelectIter(TreeIter)
	UnselectPath(TreePath)
}

func AssertTreeSelection(_ TreeSelection) {}
