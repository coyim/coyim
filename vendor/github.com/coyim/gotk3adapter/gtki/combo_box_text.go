package gtki

type ComboBoxText interface {
	ComboBox

	AppendText(string)
	GetActiveText() string
	RemoveAll()
}

func AssertComboBoxText(_ ComboBoxText) {}
