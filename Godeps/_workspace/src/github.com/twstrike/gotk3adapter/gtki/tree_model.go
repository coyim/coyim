package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type TreeModel interface {
	GetIter(TreePath) (TreeIter, error)
	GetIterFirst() (TreeIter, bool)
	GetIterFromString(string) (TreeIter, error)
	GetPath(TreeIter) (TreePath, error)
	GetValue(TreeIter, int) (glibi.Value, error)
	IterNext(TreeIter) bool
}

func AssertTreeModel(_ TreeModel) {}
