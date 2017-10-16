package glib_mock

import "github.com/coyim/gotk3adapter/glibi"

type MockObject struct {
	refCount int
}

func (*MockObject) Connect(v1 string, v2 interface{}, v3 ...interface{}) (glibi.SignalHandle, error) {
	return glibi.SignalHandle(0), nil
}

func (*MockObject) ConnectAfter(v1 string, v2 interface{}, v3 ...interface{}) (glibi.SignalHandle, error) {
	return glibi.SignalHandle(0), nil
}

func (*MockObject) Emit(v1 string, v2 ...interface{}) (interface{}, error) {
	return nil, nil
}

func (*MockObject) GetProperty(string) (interface{}, error) {
	return nil, nil
}

func (*MockObject) SetProperty(v1 string, v2 interface{}) error {
	return nil
}

func (o *MockObject) Ref() {
	o.refCount++
}

func (o *MockObject) Unref() {
	o.refCount--
}
