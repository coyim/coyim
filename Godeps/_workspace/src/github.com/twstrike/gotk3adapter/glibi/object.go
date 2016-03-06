package glibi

type Object interface {
	Connect(string, interface{}, ...interface{}) (SignalHandle, error)
	Emit(string, ...interface{}) (interface{}, error)
	SetProperty(string, interface{}) error
}

func AssertObject(_ Object) {}
