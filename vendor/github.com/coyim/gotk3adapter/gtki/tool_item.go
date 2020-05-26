package gtki

type ToolItem interface {
	Bin
}

func AssertToolItem(_ ToolItem) {}
