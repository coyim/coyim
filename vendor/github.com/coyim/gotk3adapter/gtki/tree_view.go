package gtki

type TreeView interface {
	Container

	CollapseRow(TreePath) bool
	ExpandAll()
	GetCursor() (TreePath, TreeViewColumn)
	GetPathAtPos(int, int, TreePath, TreeViewColumn, *int, *int) bool
	GetSelection() (TreeSelection, error)
	SetEnableSearch(bool)
	GetEnableSearch() bool
	SetSearchColumn(int)
	GetSearchColumn() int
	SetSearchEntry(Entry)
	GetSearchEntry() Entry
	SetSearchEqualSubstringMatch()
	GetModel() (TreeModel, error)
	SetModel(TreeModel)
}

func AssertTreeView(_ TreeView) {}
