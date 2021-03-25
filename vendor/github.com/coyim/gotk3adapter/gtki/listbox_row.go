package gtki

type ListBoxRow interface {
	Bin
	GetIndex() int
}

func AssertListBoxRow(_ ListBoxRow) {}
