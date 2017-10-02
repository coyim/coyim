package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockSpinButton struct {
	MockEntry
}

func (*MockSpinButton) GetValueAsInt() int {
	return 0
}

func (*MockSpinButton) SetValue(float64) {
}

func (*MockSpinButton) GetValue() float64 {
	return 0
}

func (*MockSpinButton) GetAdjustment() gtki.Adjustment {
	return nil
}

func (*MockSpinButton) SetRange(float64, float64) {
}

func (*MockSpinButton) SetIncrements(float64, float64) {
}
