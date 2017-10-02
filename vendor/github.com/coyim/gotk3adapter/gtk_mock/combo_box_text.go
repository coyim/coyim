package gtk_mock

type MockComboBoxText struct {
	MockComboBox
}

func (*MockComboBoxText) AppendText(v1 string) {
}

func (*MockComboBoxText) GetActiveText() string {
	return ""
}
