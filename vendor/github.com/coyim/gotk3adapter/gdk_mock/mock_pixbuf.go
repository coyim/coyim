package gdk_mock

type MockPixbuf struct {
}

func (*MockPixbuf) SavePNG(string, int) error {
	return nil
}
