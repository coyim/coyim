package gtki

import "github.com/coyim/gotk3adapter/glibi"

type CellRenderer interface {
	glibi.Object
}

func AssertCellRenderer(_ CellRenderer) {}
