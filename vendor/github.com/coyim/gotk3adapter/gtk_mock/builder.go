package gtk_mock

import (
	"github.com/coyim/gotk3adapter/glib_mock"
	"github.com/coyim/gotk3adapter/glibi"
)

type MockBuilder struct {
	glib_mock.MockObject
}

func (*MockBuilder) AddFromString(v1 string) error {
	return nil
}

func (*MockBuilder) AddFromResource(v1 string) error {
	return nil
}

func (*MockBuilder) GetObject(v1 string) (glibi.Object, error) {
	return nil, nil
}

func (*MockBuilder) ConnectSignals(v1 map[string]interface{}) {
}
