package gtki

type Revealer interface {
	Bin

	SetRevealChild(bool)
	GetRevealChild() bool
}

func AssertRevealer(_ Revealer) {}
