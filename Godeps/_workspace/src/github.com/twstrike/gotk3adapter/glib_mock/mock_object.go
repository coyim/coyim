package glib_mock

import "github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"

type MockObject struct{}

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
