package gtki

type FileChooserDialog interface {
	Dialog
	FileChooser
}

func AssertFileChooserDialog(_ FileChooserDialog) {}
