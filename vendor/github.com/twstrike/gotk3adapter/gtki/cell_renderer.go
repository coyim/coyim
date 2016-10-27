package gtki

import "github.com/twstrike/gotk3adapter/glibi"

type CellRenderer interface {
	glibi.Object
}

func AssertCellRenderer(_ CellRenderer) {}
