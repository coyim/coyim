package gtki

type TreePath interface {
	GetDepth() int
	String() string
}

func AssertTreePath(_ TreePath) {}
