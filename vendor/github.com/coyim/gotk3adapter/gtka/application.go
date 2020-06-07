package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type application struct {
	*gliba.Application
	internal *gtk.Application
}

func WrapApplicationSimple(v *gtk.Application) gtki.Application {
	if v == nil {
		return nil
	}
	return &application{gliba.WrapApplicationSimple(&v.Application), v}
}

func WrapApplication(v *gtk.Application, e error) (gtki.Application, error) {
	return WrapApplicationSimple(v), e
}

func UnwrapApplication(v gtki.Application) *gtk.Application {
	if v == nil {
		return nil
	}
	return v.(*application).internal
}

func (v *application) GetActiveWindow() gtki.Window {
	ret := WrapWindowSimple(v.internal.GetActiveWindow())
	if ret == nil {
		return nil
	}
	return ret
}

func (v *application) AddWindow(w gtki.Window) {
	v.internal.AddWindow(UnwrapWindow(w))
}

func (v *application) RemoveWindow(w gtki.Window) {
	v.internal.RemoveWindow(UnwrapWindow(w))
}

func (v *application) PrefersAppMenu() bool {
	return v.internal.PrefersAppMenu()
}

func (v *application) GetAppMenu() glibi.MenuModel {
	return gliba.WrapMenuModelSimple(v.internal.GetAppMenu())
}

func (v *application) SetAppMenu(val glibi.MenuModel) {
	v.internal.SetAppMenu(gliba.UnwrapMenuModel(val))
}

func (v *application) GetMenubar() glibi.MenuModel {
	return gliba.WrapMenuModelSimple(v.internal.GetMenubar())
}

func (v *application) SetMenubar(val glibi.MenuModel) {
	v.internal.SetMenubar(gliba.UnwrapMenuModel(val))
}

func (v *application) LookupAction(actionName string) glibi.Action {
	return gliba.WrapAction(v.internal.LookupAction(actionName))
}

func (v *application) AddAction(action glibi.Action) {
	v.internal.AddAction(gliba.UnwrapAction(action))
}

func (v *application) RemoveAction(actionName string) {
	v.internal.RemoveAction(actionName)
}

func (v *application) HasAction(actionName string) bool {
	return v.internal.HasAction(actionName)
}

func (v *application) GetActionEnabled(actionName string) bool {
	return v.internal.GetActionEnabled(actionName)
}

func (v *application) GetActionParameterType(actionName string) glibi.VariantType {
	return gliba.WrapVariantType(v.internal.GetActionParameterType(actionName))
}

func (v *application) GetActionStateType(actionName string) glibi.VariantType {
	return gliba.WrapVariantType(v.internal.GetActionStateType(actionName))
}

func (v *application) GetActionState(actionName string) glibi.Variant {
	return gliba.WrapVariant(v.internal.GetActionState(actionName))
}

func (v *application) GetActionStateHint(actionName string) glibi.Variant {
	return gliba.WrapVariant(v.internal.GetActionStateHint(actionName))
}

func (v *application) ChangeActionState(actionName string, value glibi.Variant) {
	v.internal.ChangeActionState(actionName, gliba.UnwrapVariant(value))
}

func (v *application) Activate(actionName string, parameter glibi.Variant) {
	v.internal.IActionGroup.Activate(actionName, gliba.UnwrapVariant(parameter))
}
