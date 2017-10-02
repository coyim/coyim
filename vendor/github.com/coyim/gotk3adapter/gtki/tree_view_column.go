package gtki

import "github.com/coyim/gotk3adapter/glibi"

type TreeViewColumn interface {
	glibi.Object
}

func AssertTreeViewColumn(_ TreeViewColumn) {}
