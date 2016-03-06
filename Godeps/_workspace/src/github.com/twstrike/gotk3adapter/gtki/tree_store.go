package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type TreeStore interface {
	glibi.Object
	TreeModel

	Append(TreeIter) TreeIter
	Clear()
	SetValue(TreeIter, int, interface{}) error
}

func AssertTreeStore(_ TreeStore) {}
