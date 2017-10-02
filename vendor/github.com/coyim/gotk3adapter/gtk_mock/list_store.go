package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockListStore struct {
	glib_mock.MockObject
}

func (*MockListStore) Clear() {
}

func (*MockListStore) Append() gtki.TreeIter {
	return nil
}

func (*MockListStore) Remove(v2 gtki.TreeIter) bool {
	return false
}

func (*MockListStore) Set2(v2 gtki.TreeIter, v3 []int, v4 []interface{}) error {
	return nil
}

func (*MockListStore) SetValue(v1 gtki.TreeIter, v2 int, v3 interface{}) error {
	return nil
}

func (*MockListStore) GetIter(v1 gtki.TreePath) (gtki.TreeIter, error) {
	return nil, nil
}

func (*MockListStore) GetIterFirst() (gtki.TreeIter, bool) {
	return nil, false
}

func (*MockListStore) GetIterFromString(v1 string) (gtki.TreeIter, error) {
	return nil, nil
}

func (*MockListStore) GetPath(v1 gtki.TreeIter) (gtki.TreePath, error) {
	return nil, nil
}

func (*MockListStore) GetValue(v1 gtki.TreeIter, v2 int) (glibi.Value, error) {
	return nil, nil
}

func (*MockListStore) IterNext(v1 gtki.TreeIter) bool {
	return false
}
