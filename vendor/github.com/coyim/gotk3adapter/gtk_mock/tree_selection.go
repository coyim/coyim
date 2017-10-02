package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockTreeSelection struct {
	glib_mock.MockObject
}

func (*MockTreeSelection) SelectIter(v1 gtki.TreeIter) {
}

func (*MockTreeSelection) UnselectPath(v1 gtki.TreePath) {
}

func (*MockTreeSelection) GetSelected() (gtki.TreeModel, gtki.TreeIter, bool) {
	return nil, nil, false
}
