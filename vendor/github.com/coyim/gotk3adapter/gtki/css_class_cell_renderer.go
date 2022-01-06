package gtki

type CSSClassCellRenderer interface {
	CellRenderer

	SetReal(CellRenderer)
}

func AssertCSSClassCellRenderer(_ CSSClassCellRenderer) {}
