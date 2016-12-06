package gtk_mock

import "github.com/twstrike/gotk3adapter/gtki"

type MockAssistant struct {
	MockWindow
}

func (a *MockAssistant) Commit() {
}

func (a *MockAssistant) NextPage() {
}

func (a *MockAssistant) PreviousPage() {
}

func (a *MockAssistant) AppendPage(page gtki.Widget) int {
	return 0
}

func (a *MockAssistant) SetPageComplete(page gtki.Widget, complete bool) {
}

func (a *MockAssistant) GetPageComplete(page gtki.Widget) bool {
	return true
}

func (a *MockAssistant) GetCurrentPage() int {
	return 0
}

func (a *MockAssistant) GetNthPage(pageNum int) (gtki.Widget, error) {
	return nil, nil
}
