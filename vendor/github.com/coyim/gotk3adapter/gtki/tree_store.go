package gtki

import "github.com/coyim/gotk3adapter/glibi"

type TreeStore interface {
	glibi.Object
	TreeModel

	Append(TreeIter) TreeIter
	Clear()
	SetValue(TreeIter, int, interface{}) error
	Remove(TreeIter) bool
}

func AssertTreeStore(_ TreeStore) {}
