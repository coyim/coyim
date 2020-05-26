package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockTreeView struct {
	MockContainer
}

func (*MockTreeView) CollapseRow(v1 gtki.TreePath) bool {
	return false
}

func (*MockTreeView) ExpandAll() {
}

func (*MockTreeView) GetCursor() (gtki.TreePath, gtki.TreeViewColumn) {
	return nil, nil
}

func (*MockTreeView) GetSelection() (gtki.TreeSelection, error) {
	return nil, nil
}

func (*MockTreeView) GetPathAtPos(v1 int, v2 int) (gtki.TreePath, gtki.TreeViewColumn, int, int, bool) {
	return nil, nil, 0, 0, false
}

func (*MockTreeView) SetEnableSearch(v1 bool) {
}

func (*MockTreeView) GetEnableSearch() bool {
	return false
}

func (*MockTreeView) SetSearchColumn(v1 int) {
}

func (*MockTreeView) GetSearchColumn() int {
	return 0
}

func (*MockTreeView) SetSearchEntry(v1 gtki.Entry) {
}

func (*MockTreeView) GetSearchEntry() gtki.Entry {
	return nil
}

func (*MockTreeView) SetSearchEqualSubstringMatch() {

}

func (*MockTreeView) GetModel() (gtki.TreeModel, error) {
	return nil, nil
}

func (*MockTreeView) SetModel(gtki.TreeModel) {
}
