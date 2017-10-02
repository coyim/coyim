package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockTreeStore struct {
	glib_mock.MockObject
}

func (v *MockTreeStore) GetIter(v1 gtki.TreePath) (gtki.TreeIter, error) {
	return nil, nil
}

func (v *MockTreeStore) GetIterFirst() (gtki.TreeIter, bool) {
	return nil, false
}

func (v *MockTreeStore) GetIterFromString(v1 string) (gtki.TreeIter, error) {
	return nil, nil
}

func (v *MockTreeStore) GetPath(v1 gtki.TreeIter) (gtki.TreePath, error) {
	return nil, nil
}

func (v *MockTreeStore) GetValue(v1 gtki.TreeIter, v2 int) (glibi.Value, error) {
	return nil, nil
}

func (v *MockTreeStore) IterNext(v1 gtki.TreeIter) bool {
	return false
}

func (v *MockTreeStore) Append(v1 gtki.TreeIter) gtki.TreeIter {
	return nil
}

func (v *MockTreeStore) Clear() {
}

func (v *MockTreeStore) SetValue(v1 gtki.TreeIter, v2 int, v3 interface{}) error {
	return nil
}
