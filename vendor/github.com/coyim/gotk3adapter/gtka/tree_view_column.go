package gtka

import "github.com/coyim/gotk3adapter/gliba"
import "github.com/coyim/gotk3adapter/gtki"
import "github.com/gotk3/gotk3/gtk"

type treeViewColumn struct {
	*gliba.Object
	internal *gtk.TreeViewColumn
}

func wrapTreeViewColumnSimple(v *gtk.TreeViewColumn) *treeViewColumn {
	if v == nil {
		return nil
	}
	return &treeViewColumn{gliba.WrapObjectSimple(v.Object), v}
}

func wrapTreeViewColumn(v *gtk.TreeViewColumn, e error) (*treeViewColumn, error) {
	return wrapTreeViewColumnSimple(v), e
}

func unwrapTreeViewColumn(v gtki.TreeViewColumn) *gtk.TreeViewColumn {
	if v == nil {
		return nil
	}
	return v.(*treeViewColumn).internal
}
