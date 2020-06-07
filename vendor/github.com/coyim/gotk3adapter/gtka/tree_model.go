package gtka

import "github.com/coyim/gotk3adapter/gtki"
import "github.com/gotk3/gotk3/gtk"

func UnwrapTreeModel(s gtki.TreeModel) gtk.ITreeModel {
	if s == nil {
		return nil
	}

	switch ss := s.(type) {
	case *listStore:
		return UnwrapListStore(ss)
	case *treeStore:
		return UnwrapTreeStore(ss)
	}
	return nil
}

func WrapTreeModelSimple(s gtk.ITreeModel) gtki.TreeModel {
	if s == nil {
		return nil
	}

	switch ss := s.(type) {
	case *gtk.ListStore:
		return WrapListStoreSimple(ss)
	case *gtk.TreeStore:
		return WrapTreeStoreSimple(ss)
	}
	return nil
}
