package gtki

type ComboBox interface {
	Bin
	CellLayout

	GetActive() int
	GetActiveIter() (TreeIter, error)
	GetActiveID() string
	SetActive(int)
	SetModel(TreeModel)
	SetIDColumn(int)
	SetEntryTextColumn(int)
	GetToggleButton() (Button, error)
}

func AssertComboBox(_ ComboBox) {}
