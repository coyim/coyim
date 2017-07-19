package gtk_mock

// MockProgressBar is a mock of the representation of GTK's GtkProgressBar.
type MockProgressBar struct {
	MockWidget
}

func (*MockProgressBar) SetFraction(float64) {
}

func (*MockProgressBar) GetFraction() float64 {
	return 0
}

func (*MockProgressBar) SetShowText(bool) {
}

func (*MockProgressBar) GetShowText() bool {
	return false
}

func (*MockProgressBar) SetText(string) {
}
