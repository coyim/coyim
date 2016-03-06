package gtki

type Editable interface {
	SetEditable(bool)
}

func AssertEditable(_ Editable) {}
