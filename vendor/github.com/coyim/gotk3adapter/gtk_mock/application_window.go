package gtk_mock

import "github.com/coyim/gotk3adapter/glibi"

type MockApplicationWindow struct {
	MockWindow
}

func (*MockApplicationWindow) SetShowMenubar(bool) {}
func (*MockApplicationWindow) GetShowMenubar() bool {
	return false
}
func (*MockApplicationWindow) GetID() uint {
	return 0
}

func (*MockApplicationWindow) LookupAction(actionName string) glibi.Action {
	return nil
}

func (*MockApplicationWindow) AddAction(action glibi.Action) {
}

func (*MockApplicationWindow) RemoveAction(actionName string) {
}

func (*MockApplicationWindow) HasAction(actionName string) bool {
	return false
}

func (*MockApplicationWindow) GetActionEnabled(actionName string) bool {
	return false
}

func (*MockApplicationWindow) GetActionParameterType(actionName string) glibi.VariantType {
	return nil
}

func (*MockApplicationWindow) GetActionStateType(actionName string) glibi.VariantType {
	return nil
}

func (*MockApplicationWindow) GetActionState(actionName string) glibi.Variant {
	return nil
}

func (*MockApplicationWindow) GetActionStateHint(actionName string) glibi.Variant {
	return nil
}

func (*MockApplicationWindow) ChangeActionState(actionName string, value glibi.Variant) {
}

func (*MockApplicationWindow) Activate(actionName string, parameter glibi.Variant) {
}
