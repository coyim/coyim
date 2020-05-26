package gtki

type FileChooser interface {
	GetFilename() string
	SetCurrentName(string)
	SetDoOverwriteConfirmation(bool)
}

func AssertFileChooser(_ FileChooser) {}
