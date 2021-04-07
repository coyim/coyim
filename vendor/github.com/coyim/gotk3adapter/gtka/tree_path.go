package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type treePath struct {
	*gtk.TreePath
}

func WrapTreePathSimple(v *gtk.TreePath) gtki.TreePath {
	if v == nil {
		return nil
	}
	return &treePath{v}
}

func WrapTreePath(v *gtk.TreePath, e error) (gtki.TreePath, error) {
	return WrapTreePathSimple(v), e
}

func UnwrapTreePath(v gtki.TreePath) *gtk.TreePath {
	if v == nil {
		return nil
	}
	return v.(*treePath).TreePath
}

func (v *treePath) GetDepth() int {
	return v.TreePath.GetDepth()
}
