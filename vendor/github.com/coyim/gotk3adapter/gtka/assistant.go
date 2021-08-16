package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3extra"
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
	obj, e := a.internal.GetNthPage(pageNum)
	return Wrap(obj).(gtki.Widget), e
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

func (a *assistant) SetPageTitle(page gtki.Widget, title string) {
	a.internal.SetPageTitle(UnwrapWidget(page), title)
}

const (
	assistantButtonBackCancelName = "cancel"
	assistantButtonBackLastName   = "back"
	assistantButtonLastName       = "last"
	assistantButtonForwardName    = "forward"
	assistantButtonApplyName      = "apply"
)

var assistantButtons = []string{
	assistantButtonBackCancelName,
	assistantButtonBackLastName,
	assistantButtonLastName,
	assistantButtonForwardName,
	assistantButtonApplyName,
}

func (a *assistant) GetButtons() []gtki.Button {
	buttons := []gtki.Button{}

	for _, buttonName := range assistantButtons {
		if obj, err := a.TemplateChild(buttonName); err == nil {
			b := WrapButtonSimple(gotk3extra.WrapButton(gliba.UnwrapObject(obj)))
			buttons = append(buttons, b)
		}
	}

	return buttons
}

func (a *assistant) GetButtonSizeGroup() (gtki.SizeGroup, error) {
	v, err := gotk3extra.GetAssistantButtonSizeGroup(a.internal)
	return WrapSizeGroup(v, err)
}

const (
	assistantActionAreaName     = "action_area"
	assistantHeaderBarName      = "headerbar"
	assistantSidebarName        = "sidebar"
	assistantContentWrapperName = "content_box"
	assistantContentName        = "content"
)

func (a *assistant) GetHeaderBar() (gtki.HeaderBar, error) {
	obj, err := a.TemplateChild(assistantHeaderBarName)
	return WrapHeaderBarSimple(gotk3extra.WrapHeaderBar(gliba.UnwrapObject(obj))), err
}

func (a *assistant) GetSidebar() (gtki.Box, error) {
	obj, err := a.TemplateChild(assistantSidebarName)
	return WrapBoxSimple(gotk3extra.WrapBox(gliba.UnwrapObject(obj))), err
}

func (a *assistant) GetNotebook() (gtki.Notebook, error) {
	obj, err := a.TemplateChild(assistantContentName)
	return WrapNotebookSimple(gotk3extra.WrapNotebook(gliba.UnwrapObject(obj))), err
}

func (a *assistant) HideBottomActionArea() {
	if obj, err := a.TemplateChild(assistantActionAreaName); err == nil {
		box := WrapBoxSimple(gotk3extra.WrapBox(gliba.UnwrapObject(obj)))
		box.Hide()
	}
}
