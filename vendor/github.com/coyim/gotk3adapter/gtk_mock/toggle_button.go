package gtk_mock

type MockToggleButton struct {
	MockButton
}

func (*MockToggleButton) GetActive() bool {
	return false
}

func (*MockToggleButton) SetActive(bool) {
}
