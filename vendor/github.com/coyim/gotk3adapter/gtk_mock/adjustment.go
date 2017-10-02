package gtk_mock

import "github.com/coyim/gotk3adapter/glib_mock"

type MockAdjustment struct {
	glib_mock.MockObject
}

func (*MockAdjustment) GetLower() float64 {
	return 0
}

func (*MockAdjustment) GetPageSize() float64 {
	return 0
}

func (*MockAdjustment) GetUpper() float64 {
	return 0
}

func (*MockAdjustment) SetValue(v1 float64) {
}
