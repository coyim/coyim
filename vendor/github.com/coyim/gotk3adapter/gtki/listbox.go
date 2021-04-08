package gtki

type ListBox interface {
	Container

	SelectRow(ListBoxRow)
	GetRowAtIndex(int) ListBoxRow
	GetSelectedRow() ListBoxRow
}

func AssertListBox(_ ListBox) {}
