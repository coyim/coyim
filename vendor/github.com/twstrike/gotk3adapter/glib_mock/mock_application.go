package glib_mock

type MockApplication struct {
	MockObject
}

func (*MockApplication) Quit() {}
func (*MockApplication) Run([]string) int {
	return 0
}
