package gtki

type ListBox interface {
	Container

	SelectRow(ListBoxRow)
	GetRowAtIndex(int) ListBoxRow
}

func AssertListBox(_ ListBox) {}
