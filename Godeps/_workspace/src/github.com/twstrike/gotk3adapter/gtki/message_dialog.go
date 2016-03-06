package gtki

type MessageDialog interface {
	Dialog
}

func AssertMessageDialog(_ MessageDialog) {}
