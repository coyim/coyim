package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

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
