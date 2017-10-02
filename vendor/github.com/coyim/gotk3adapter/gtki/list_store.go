package gtki

import "github.com/coyim/gotk3adapter/glibi"

type ListStore interface {
	glibi.Object
	TreeModel

	Append() TreeIter
	Clear()
	Remove(TreeIter) bool
	Set2(TreeIter, []int, []interface{}) error
	SetValue(TreeIter, int, interface{}) error
}

func AssertListStore(_ ListStore) {}
