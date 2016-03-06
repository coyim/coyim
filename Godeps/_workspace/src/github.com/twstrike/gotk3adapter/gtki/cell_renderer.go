package gtki

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type CellRenderer interface {
	glibi.Object
}

func AssertCellRenderer(_ CellRenderer) {}
