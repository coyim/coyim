package gtki

// ProgressBar is an interface of Gtk.ProgressBar
type ProgressBar interface {
	Widget

	SetFraction(float64)
	GetFraction() float64
	SetShowText(bool)
	GetShowText() bool
	SetText(string)
}

// AssertProgressBar asserts the ProgressBar
func AssertProgressBar(_ ProgressBar) {}
