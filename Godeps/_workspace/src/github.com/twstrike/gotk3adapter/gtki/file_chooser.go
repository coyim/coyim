package gtki

type FileChooser interface {
	GetFilename() string
	SetCurrentName(string)
}

func AssertFileChooser(_ FileChooser) {}
