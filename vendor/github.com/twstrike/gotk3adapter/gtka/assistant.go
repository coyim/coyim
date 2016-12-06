package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/gotk3adapter/gtki"
)

type assistant struct {
	*window
	internal *gtk.Assistant
}

func wrapAssistantSimple(v *gtk.Assistant) *assistant {
	if v == nil {
		return nil
	}
	return &assistant{wrapWindowSimple(&v.Window), v}
}

func wrapAssistant(v *gtk.Assistant, e error) (*assistant, error) {
	return wrapAssistantSimple(v), e
}

func unwrapAssistant(v gtki.Assistant) *gtk.Assistant {
	if v == nil {
		return nil
	}
	return v.(*assistant).internal
}

func (a *assistant) Commit() {
	a.internal.Commit()
}

func (a *assistant) NextPage() {
	a.internal.NextPage()
}

func (a *assistant) PreviousPage() {
	a.internal.PreviousPage()
}

func (a *assistant) AppendPage(page gtki.Widget) int {
	return a.internal.AppendPage(unwrapWidget(page))
}

func (a *assistant) SetPageComplete(page gtki.Widget, complete bool) {
	a.internal.SetPageComplete(unwrapWidget(page), complete)
}

func (a *assistant) GetPageComplete(page gtki.Widget) bool {
	return a.internal.GetPageComplete(unwrapWidget(page))
}

func (a *assistant) GetCurrentPage() int {
	return a.internal.GetCurrentPage()
}

func (a *assistant) GetNthPage(pageNum int) (gtki.Widget, error) {
	return wrapWidget(a.internal.GetNthPage(pageNum))
}
