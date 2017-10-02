package gtki

type ComboBoxText interface {
	ComboBox

	AppendText(string)
	GetActiveText() string
}

func AssertComboBoxText(_ ComboBoxText) {}
