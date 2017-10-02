package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockNotebook struct {
	MockContainer
}

func (*MockNotebook) NextPage() {
}

func (*MockNotebook) PrevPage() {
}

func (*MockNotebook) GetCurrentPage() int {
	return 0
}

func (*MockNotebook) GetNPages() int {
	return 0
}

func (*MockNotebook) SetCurrentPage(v1 int) {
}

func (*MockNotebook) SetShowTabs(v1 bool) {
}

func (*MockNotebook) AppendPage(v1, v2 gtki.Widget) int {
	return 0
}

func (*MockNotebook) GetNthPage(v1 int) (gtki.Widget, error) {
	return nil, nil
}

func (*MockNotebook) SetTabLabelText(v1 gtki.Widget, v2 string) {
}
