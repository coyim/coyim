package glibi

type Action interface {
	Object

	GetName() string
	GetEnabled() bool
	GetState() Variant
	GetStateHint() Variant
	GetParameterType() VariantType
	GetStateType() VariantType
	ChangeState(value Variant)
	Activate(parameter Variant)
}

func AssertAction(_ Action) {}

type SimpleAction interface {
	Action

	SetEnabled(enabled bool)
	SetState(value Variant)
}

func AssertSimpleAction(_ SimpleAction) {}

type PropertyAction interface {
	Action
}

func AssertPropertyAction(_ PropertyAction) {}
