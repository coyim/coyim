package gtki

// Spinner is an interface of Gtk.Spinner
type Spinner interface {
	Widget

	Start()
	Stop()
}

// AssertSpinner asserts the spinner
func AssertSpinner(_ Spinner) {}
