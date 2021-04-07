package gtki

type TreePath interface {
	GetDepth() int
}

func AssertTreePath(_ TreePath) {}
