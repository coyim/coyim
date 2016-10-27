package gtki

import "github.com/twstrike/gotk3adapter/glibi"

type TreeViewColumn interface {
	glibi.Object
}

func AssertTreeViewColumn(_ TreeViewColumn) {}
