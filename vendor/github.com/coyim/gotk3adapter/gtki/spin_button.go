package gtki

type SpinButton interface {
	Entry

	GetValueAsInt() int
	SetValue(float64)
	GetValue() float64
	GetAdjustment() Adjustment
	SetRange(float64, float64)
	SetIncrements(float64, float64)
}

func AssertSpinButton(_ SpinButton) {}
