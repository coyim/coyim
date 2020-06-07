package gtk_mock

type MockApplicationWindow struct {
	MockWindow
}

func (*MockApplicationWindow) SetShowMenubar(bool) {}
func (*MockApplicationWindow) GetShowMenubar() bool {
	return false
}
func (*MockApplicationWindow) GetID() uint {
	return 0
}
