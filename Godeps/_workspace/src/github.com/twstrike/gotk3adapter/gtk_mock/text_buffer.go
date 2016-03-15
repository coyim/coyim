package gtk_mock

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
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

func (*MockTextBuffer) GetStartIter() gtki.TextIter {
	return nil
}

func (*MockTextBuffer) Insert(v1 gtki.TextIter, v2 string) {
}

func (*MockTextBuffer) GetText(gtki.TextIter, gtki.TextIter, bool) string {
	return ""
}
