package pango_mock

type MockFontDescription struct{}

func (*MockFontDescription) GetSize() int {
	return 0
}
