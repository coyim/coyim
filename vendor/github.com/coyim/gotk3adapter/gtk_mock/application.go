package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type MockApplication struct {
	glib_mock.MockApplication
}

func (*MockApplication) GetActiveWindow() gtki.Window {
	return nil
}

func (*MockApplication) AddWindow(gtki.Window)    {}
func (*MockApplication) RemoveWindow(gtki.Window) {}
func (*MockApplication) PrefersAppMenu() bool {
	return false
}

func (*MockApplication) GetAppMenu() glibi.MenuModel {
	return nil
}

func (*MockApplication) SetAppMenu(glibi.MenuModel) {
}

func (*MockApplication) GetMenubar() glibi.MenuModel {
	return nil
}

func (*MockApplication) SetMenubar(glibi.MenuModel) {
}

func (*MockApplication) LookupAction(actionName string) glibi.Action {
	return nil
}

func (*MockApplication) AddAction(action glibi.Action) {
}

func (*MockApplication) RemoveAction(actionName string) {
}

func (*MockApplication) HasAction(actionName string) bool {
	return false
}

func (*MockApplication) GetActionEnabled(actionName string) bool {
	return false
}

func (*MockApplication) GetActionParameterType(actionName string) glibi.VariantType {
	return nil
}

func (*MockApplication) GetActionStateType(actionName string) glibi.VariantType {
	return nil
}

func (*MockApplication) GetActionState(actionName string) glibi.Variant {
	return nil
}

func (*MockApplication) GetActionStateHint(actionName string) glibi.Variant {
	return nil
}

func (*MockApplication) ChangeActionState(actionName string, value glibi.Variant) {
}

func (*MockApplication) Activate(actionName string, parameter glibi.Variant) {
}
