package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type action struct {
	*Object
	*glib.Action
}

func WrapAction(v *glib.Action) glibi.Action {
	if v == nil {
		return nil
	}
	return &action{WrapObjectSimple(v.Object), v}
}

func UnwrapAction(v glibi.Action) *glib.Action {
	if v == nil {
		return nil
	}
	return v.(*action).Action
}

func (v *action) GetName() string {
	return v.Action.GetName()
}

func (v *action) GetEnabled() bool {
	return v.Action.GetEnabled()
}

func (v *action) GetState() glibi.Variant {
	return WrapVariant(v.Action.GetState())
}

func (v *action) GetStateHint() glibi.Variant {
	return WrapVariant(v.Action.GetStateHint())
}

func (v *action) GetParameterType() glibi.VariantType {
	return WrapVariantType(v.Action.GetParameterType())
}

func (v *action) GetStateType() glibi.VariantType {
	return WrapVariantType(v.Action.GetStateType())
}

func (v *action) ChangeState(value glibi.Variant) {
	v.Action.ChangeState(UnwrapVariant(value))
}

func (v *action) Activate(parameter glibi.Variant) {
	v.Action.Activate(UnwrapVariant(parameter))
}

type simpleAction struct {
	*action
	*glib.SimpleAction
}

func WrapSimpleAction(v *glib.SimpleAction) glibi.SimpleAction {
	if v == nil {
		return nil
	}
	return &simpleAction{WrapAction(&v.Action).(*action), v}
}

func UnwrapSimpleAction(v glibi.SimpleAction) *glib.SimpleAction {
	if v == nil {
		return nil
	}
	return v.(*simpleAction).SimpleAction
}

func (v *simpleAction) SetEnabled(enabled bool) {
	v.SimpleAction.SetEnabled(enabled)
}

func (v *simpleAction) SetState(value glibi.Variant) {
	v.SimpleAction.SetState(UnwrapVariant(value))
}

type propertyAction struct {
	*action
	*glib.PropertyAction
}

func WrapPropertyAction(v *glib.PropertyAction) glibi.PropertyAction {
	if v == nil {
		return nil
	}
	return &propertyAction{WrapAction(&v.Action).(*action), v}
}

func UnwrapPropertyAction(v glibi.PropertyAction) *glib.PropertyAction {
	if v == nil {
		return nil
	}
	return v.(*propertyAction).PropertyAction
}
