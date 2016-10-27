package gtki

type ScrolledWindow interface {
	Bin

	GetVAdjustment() Adjustment
}

func AssertScrolledWindow(_ ScrolledWindow) {}
