package glibi

type ActionGroup interface {
	HasAction(actionName string) bool
	GetActionEnabled(actionName string) bool
	GetActionParameterType(actionName string) VariantType
	GetActionStateType(actionName string) VariantType
	GetActionState(actionName string) Variant
	GetActionStateHint(actionName string) Variant
	ChangeActionState(actionName string, value Variant)
	Activate(actionName string, parameter Variant)
}
