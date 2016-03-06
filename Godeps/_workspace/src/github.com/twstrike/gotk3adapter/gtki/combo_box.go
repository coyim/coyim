package gtki

type ComboBox interface {
	Bin
	CellLayout

	GetActiveIter() (TreeIter, error)
	SetActive(int)
	SetModel(TreeModel)
}

func AssertComboBox(_ ComboBox) {}
