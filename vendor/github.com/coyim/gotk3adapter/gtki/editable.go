package gtki

type Editable interface {
	SetEditable(bool)
	SetPosition(int)
	GetPosition() int
}

func AssertEditable(_ Editable) {}
