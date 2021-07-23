package glibi

type Object interface {
	Connect(string, interface{}) SignalHandle
	ConnectAfter(string, interface{}) SignalHandle
	Emit(string, ...interface{}) (interface{}, error)
	GetProperty(string) (interface{}, error)
	Ref()
	SetProperty(string, interface{}) error
	Unref()
}

func AssertObject(_ Object) {}
