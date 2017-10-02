package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gtki"
)

type notebook struct {
	*container
	internal *gtk.Notebook
}

func wrapNotebookSimple(v *gtk.Notebook) *notebook {
	if v == nil {
		return nil
	}
	return &notebook{wrapContainerSimple(&v.Container), v}
}

func wrapNotebook(v *gtk.Notebook, e error) (*notebook, error) {
	return wrapNotebookSimple(v), e
}

func unwrapNotebook(v gtki.Notebook) *gtk.Notebook {
	if v == nil {
		return nil
	}
	return v.(*notebook).internal
}

func (v *notebook) NextPage() {
	v.internal.NextPage()
}

func (v *notebook) PrevPage() {
	v.internal.PrevPage()
}

func (v *notebook) GetCurrentPage() int {
	return v.internal.GetCurrentPage()
}

func (v *notebook) GetNPages() int {
	return v.internal.GetNPages()
}

func (v *notebook) SetCurrentPage(v1 int) {
	v.internal.SetCurrentPage(v1)
}

func (v *notebook) SetShowTabs(v1 bool) {
	v.internal.SetShowTabs(v1)
}

func (v *notebook) AppendPage(v1, v2 gtki.Widget) int {
	return v.internal.AppendPage(unwrapWidget(v1), unwrapWidget(v2))
}

func (v *notebook) GetNthPage(v1 int) (gtki.Widget, error) {
	vx1, vx2 := v.internal.GetNthPage(v1)
	return Wrap(vx1).(gtki.Widget), vx2
}

func (v *notebook) SetTabLabelText(v1 gtki.Widget, v2 string) {
	v.internal.SetTabLabelText(unwrapWidget(v1), v2)
}
