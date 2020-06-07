package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type treeIter struct {
	*gtk.TreeIter
}

func WrapTreeIterSimple(v *gtk.TreeIter) gtki.TreeIter {
	if v == nil {
		return nil
	}
	return &treeIter{v}
}

func WrapTreeIter(v *gtk.TreeIter, e error) (gtki.TreeIter, error) {
	return WrapTreeIterSimple(v), e
}

func UnwrapTreeIter(v gtki.TreeIter) *gtk.TreeIter {
	if v == nil {
		return nil
	}
	return v.(*treeIter).TreeIter
}
