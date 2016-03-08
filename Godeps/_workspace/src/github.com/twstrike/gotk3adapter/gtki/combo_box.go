package gtki

type ComboBox interface {
	Bin
	CellLayout

	GetActiveIter() (TreeIter, error)
	GetActiveID() string
	SetActive(int)
	SetModel(TreeModel)
}

func AssertComboBox(_ ComboBox) {}
