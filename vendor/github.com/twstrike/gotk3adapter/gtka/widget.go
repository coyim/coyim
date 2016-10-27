package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/gotk3adapter/gdka"
	"github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/gotk3adapter/gliba"
	"github.com/twstrike/gotk3adapter/gtki"
)

type widget struct {
	*gliba.Object
	internal *gtk.Widget
}

type asWidget interface {
	toWidget() *widget
}

func (v *widget) toWidget() *widget {
	return v
}

func wrapWidgetSimple(v *gtk.Widget) *widget {
	if v == nil {
		return nil
	}
	return &widget{gliba.WrapObjectSimple(v.Object), v}
}

func wrapWidget(v *gtk.Widget, e error) (*widget, error) {
	return wrapWidgetSimple(v), e
}

func unwrapWidget(v gtki.Widget) *gtk.Widget {
	if v == nil {
		return nil
	}
	return v.(asWidget).toWidget().internal
}

func (v *widget) SetHExpand(v1 bool) {
	v.internal.SetHExpand(v1)
}

func (v *widget) SetSensitive(v1 bool) {
	v.internal.SetSensitive(v1)
}

func (v *widget) SetVisible(v1 bool) {
	v.internal.SetVisible(v1)
}

func (v *widget) SetName(v1 string) {
	v.internal.SetName(v1)
}

func (v *widget) SetMarginTop(v1 int) {
	v.internal.SetMarginTop(v1)
}

func (v *widget) SetMarginBottom(v1 int) {
	v.internal.SetMarginBottom(v1)
}

func (v *widget) SetSizeRequest(v1, v2 int) {
	v.internal.SetSizeRequest(v1, v2)
}

func (v *widget) GetAllocatedHeight() int {
	return v.internal.GetAllocatedHeight()
}

func (v *widget) GetAllocatedWidth() int {
	return v.internal.GetAllocatedWidth()
}

func (v *widget) GrabFocus() {
	v.internal.GrabFocus()
}

func (v *widget) GrabDefault() {
	v.internal.GrabDefault()
}

func (v *widget) Hide() {
	v.internal.Hide()
}

func (v *widget) HideOnDelete() {
	v.internal.HideOnDelete()
}

func (v *widget) Show() {
	v.internal.Show()
}

func (v *widget) ShowAll() {
	v.internal.ShowAll()
}

func (v *widget) GetWindow() (gdki.Window, error) {
	return gdka.WrapWindow(v.internal.GetWindow())
}

func (v *widget) GetStyleContext() (gtki.StyleContext, error) {
	return wrapStyleContext(v.internal.GetStyleContext())
}

func (v *widget) SetHAlign(v2 gtki.Align) {
	v.internal.SetHAlign(gtk.Align(v2))
}

func (v *widget) Destroy() {
	v.internal.Destroy()
}

func (v *widget) HasFocus() bool {
	return v.internal.HasFocus()
}
