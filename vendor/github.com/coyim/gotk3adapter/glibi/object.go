package glibi

type Object interface {
	Connect(string, interface{}, ...interface{}) (SignalHandle, error)
	ConnectAfter(string, interface{}, ...interface{}) (SignalHandle, error)
	Emit(string, ...interface{}) (interface{}, error)
	GetProperty(string) (interface{}, error)
	Ref()
	SetProperty(string, interface{}) error
	Unref()
}

func AssertObject(_ Object) {}
