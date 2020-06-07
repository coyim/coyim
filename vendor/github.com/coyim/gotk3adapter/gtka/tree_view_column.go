package gtka

import "github.com/coyim/gotk3adapter/gliba"
import "github.com/coyim/gotk3adapter/gtki"
import "github.com/gotk3/gotk3/gtk"

type treeViewColumn struct {
	*gliba.Object
	internal *gtk.TreeViewColumn
}

func WrapTreeViewColumnSimple(v *gtk.TreeViewColumn) gtki.TreeViewColumn {
	if v == nil {
		return nil
	}
	return &treeViewColumn{gliba.WrapObjectSimple(v.Object), v}
}

func WrapTreeViewColumn(v *gtk.TreeViewColumn, e error) (gtki.TreeViewColumn, error) {
	return WrapTreeViewColumnSimple(v), e
}

func UnwrapTreeViewColumn(v gtki.TreeViewColumn) *gtk.TreeViewColumn {
	if v == nil {
		return nil
	}
	return v.(*treeViewColumn).internal
}
