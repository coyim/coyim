package gtki

type CellLayout interface {
	AddAttribute(CellRenderer, string, int)
	PackStart(CellRenderer, bool)
}

func AssertCellLayout(_ CellLayout) {}
