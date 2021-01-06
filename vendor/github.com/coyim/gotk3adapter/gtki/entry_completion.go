package gtki

import (
	"github.com/coyim/gotk3adapter/glibi"
)

type EntryCompletion interface {
	glibi.Object

	SetModel(TreeModel)
	SetTextColumn(int)
	SetMinimumKeyLength(int)
}
