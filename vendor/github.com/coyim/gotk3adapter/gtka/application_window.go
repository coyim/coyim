package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type applicationWindow struct {
	*window
	internal *gtk.ApplicationWindow
}

func WrapApplicationWindowSimple(v *gtk.ApplicationWindow) gtki.ApplicationWindow {
	if v == nil {
		return nil
	}
	return &applicationWindow{WrapWindowSimple(&v.Window).(*window), v}
}

func WrapApplicationWindow(v *gtk.ApplicationWindow, e error) (gtki.ApplicationWindow, error) {
	return WrapApplicationWindowSimple(v), e
}

func UnwrapApplicationWindow(v gtki.ApplicationWindow) *gtk.ApplicationWindow {
	if v == nil {
		return nil
	}
	return v.(*applicationWindow).internal
}

func (v *applicationWindow) SetShowMenubar(val bool) {
	v.internal.SetShowMenubar(val)
}

func (v *applicationWindow) GetShowMenubar() bool {
	return v.internal.GetShowMenubar()
}

func (v *applicationWindow) GetID() uint {
	return v.internal.GetID()
}

func (v *applicationWindow) LookupAction(actionName string) glibi.Action {
	return gliba.WrapAction(v.internal.LookupAction(actionName))
}

func (v *applicationWindow) AddAction(action glibi.Action) {
	v.internal.AddAction(gliba.UnwrapAction(action))
}

func (v *applicationWindow) RemoveAction(actionName string) {
	v.internal.RemoveAction(actionName)
}

func (v *applicationWindow) HasAction(actionName string) bool {
	return v.internal.HasAction(actionName)
}

func (v *applicationWindow) GetActionEnabled(actionName string) bool {
	return v.internal.GetActionEnabled(actionName)
}

func (v *applicationWindow) GetActionParameterType(actionName string) glibi.VariantType {
	return gliba.WrapVariantType(v.internal.GetActionParameterType(actionName))
}

func (v *applicationWindow) GetActionStateType(actionName string) glibi.VariantType {
	return gliba.WrapVariantType(v.internal.GetActionStateType(actionName))
}

func (v *applicationWindow) GetActionState(actionName string) glibi.Variant {
	return gliba.WrapVariant(v.internal.GetActionState(actionName))
}

func (v *applicationWindow) GetActionStateHint(actionName string) glibi.Variant {
	return gliba.WrapVariant(v.internal.GetActionStateHint(actionName))
}

func (v *applicationWindow) ChangeActionState(actionName string, value glibi.Variant) {
	v.internal.ChangeActionState(actionName, gliba.UnwrapVariant(value))
}

func (v *applicationWindow) Activate(actionName string, parameter glibi.Variant) {
	v.internal.IActionGroup.Activate(actionName, gliba.UnwrapVariant(parameter))
}
