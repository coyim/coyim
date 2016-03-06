package gtki

type TreeView interface {
	Container

	CollapseRow(TreePath) bool
	ExpandAll()
	GetCursor() (TreePath, TreeViewColumn)
	GetPathAtPos(int, int, TreePath, TreeViewColumn, *int, *int) bool
	GetSelection() (TreeSelection, error)
}

func AssertTreeView(_ TreeView) {}
