package gdk_mock

type MockRgba struct {
}

func (*MockRgba) String() string {
	return ""
}

func (*MockRgba) GetRed() float64 {
	return 0
}

func (*MockRgba) GetGreen() float64 {
	return 0
}

func (*MockRgba) GetBlue() float64 {
	return 0
}

func (*MockRgba) GetAlpha() float64 {
	return 0
}

func (*MockRgba) SetRed(c float64) {
}

func (*MockRgba) SetGreen(c float64) {
}

func (*MockRgba) SetBlue(c float64) {
}

func (*MockRgba) SetAlpha(c float64) {
}

func (*MockRgba) Colors() (r, g, b, a float64) {
	return 0, 0, 0, 0
}

func (*MockRgba) SetColors(r, g, b, a float64) {
}

func (*MockRgba) Parse(spec string) bool {
	return false
}
