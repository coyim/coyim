package gtka

import (
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type assistant struct {
	*window
	internal *gtk.Assistant
}

func WrapAssistantSimple(v *gtk.Assistant) gtki.Assistant {
	if v == nil {
		return nil
	}
	return &assistant{WrapWindowSimple(&v.Window).(*window), v}
}

func WrapAssistant(v *gtk.Assistant, e error) (gtki.Assistant, error) {
	return WrapAssistantSimple(v), e
}

func UnwrapAssistant(v gtki.Assistant) *gtk.Assistant {
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

func (a *assistant) SetCurrentPage(pageNum int) {
	a.internal.SetCurrentPage(pageNum)
}

func (a *assistant) GetCurrentPage() int {
	return a.internal.GetCurrentPage()
}

func (a *assistant) GetNthPage(pageNum int) (gtki.Widget, error) {
	return WrapWidget(a.internal.GetNthPage(pageNum))
}

func (a *assistant) AppendPage(page gtki.Widget) int {
	return a.internal.AppendPage(UnwrapWidget(page))
}

func (a *assistant) SetPageType(page gtki.Widget, ptype gtki.AssistantPageType) {
	a.internal.SetPageType(UnwrapWidget(page), gtk.AssistantPageType(ptype))
}

func (a *assistant) GetPageType(page gtki.Widget) gtki.AssistantPageType {
	ptype := a.internal.GetPageType(UnwrapWidget(page))
	return gtki.AssistantPageType(ptype)
}

func (a *assistant) SetPageComplete(page gtki.Widget, complete bool) {
	a.internal.SetPageComplete(UnwrapWidget(page), complete)
}

func (a *assistant) GetPageComplete(page gtki.Widget) bool {
	return a.internal.GetPageComplete(UnwrapWidget(page))
}

func (a *assistant) AddActionWidget(child gtki.Widget) {
	a.internal.AddActionWidget(UnwrapWidget(child))
}

func (a *assistant) RemoveActionWidget(child gtki.Widget) {
	a.internal.RemoveActionWidget(UnwrapWidget(child))
}

func (a *assistant) UpdateButtonsState() {
	a.internal.UpdateButtonsState()
}
