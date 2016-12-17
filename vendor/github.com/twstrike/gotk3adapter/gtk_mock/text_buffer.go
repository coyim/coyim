package gtk_mock

import (
	"github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/gotk3adapter/gtki"
)

type MockTextBuffer struct {
	glib_mock.MockObject
}

func (*MockTextBuffer) ApplyTagByName(v1 string, v2, v3 gtki.TextIter) {
}

func (*MockTextBuffer) GetCharCount() int {
	return 0
}

func (*MockTextBuffer) GetEndIter() gtki.TextIter {
	return nil
}

func (*MockTextBuffer) GetIterAtOffset(v1 int) gtki.TextIter {
	return nil
}

func (*MockTextBuffer) GetLineCount() int {
	return 0
}

func (*MockTextBuffer) GetStartIter() gtki.TextIter {
	return nil
}

func (*MockTextBuffer) Insert(v1 gtki.TextIter, v2 string) {
}

func (*MockTextBuffer) InsertAtCursor(v1 string) {
}

func (*MockTextBuffer) GetText(gtki.TextIter, gtki.TextIter, bool) string {
	return ""
}

func (*MockTextBuffer) Delete(gtki.TextIter, gtki.TextIter) {
}

func (*MockTextBuffer) CreateMark(string, gtki.TextIter, bool) gtki.TextMark {
	return nil
}

func (*MockTextBuffer) GetIterAtMark(gtki.TextMark) gtki.TextIter {
	return nil
}
