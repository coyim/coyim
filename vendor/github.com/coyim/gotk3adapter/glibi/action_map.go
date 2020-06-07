package glibi

type ActionMap interface {
	LookupAction(actionName string) Action
	AddAction(action Action)
	RemoveAction(actionName string)
}
