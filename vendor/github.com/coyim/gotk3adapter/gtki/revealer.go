package gtki

type Revealer interface {
	Bin

	SetRevealChild(bool)
}

func AssertRevealer(_ Revealer) {}
