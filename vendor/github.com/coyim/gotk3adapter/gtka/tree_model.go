package gtka

import "github.com/coyim/gotk3adapter/gtki"
import "github.com/gotk3/gotk3/gtk"

func unwrapTreeModel(s gtki.TreeModel) gtk.ITreeModel {
	if s == nil {
		return nil
	}

	switch ss := s.(type) {
	case *listStore:
		return unwrapListStore(ss)
	case *treeStore:
		return unwrapTreeStore(ss)
	}
	return nil
}

func wrapTreeModelSimple(s gtk.ITreeModel) gtki.TreeModel {
	if s == nil {
		return nil
	}

	switch ss := s.(type) {
	case *gtk.ListStore:
		return wrapListStoreSimple(ss)
	case *gtk.TreeStore:
		return wrapTreeStoreSimple(ss)
	}
	return nil
}
