package gtk_mock

type MockCheckMenuItem struct {
	MockMenuItem
}

func (*MockCheckMenuItem) GetActive() bool {
	return false
}

func (*MockCheckMenuItem) SetActive(v1 bool) {
}
