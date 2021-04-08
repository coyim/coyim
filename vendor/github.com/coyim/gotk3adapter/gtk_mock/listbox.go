package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockListBox struct {
	MockContainer
}

func (*MockListBox) SelectRow(gtki.ListBoxRow) {
}

func (*MockListBox) GetRowAtIndex(int) gtki.ListBoxRow {
	return nil
}

func (*MockListBox) GetSelectedRow() gtki.ListBoxRow {
	return nul
}
