package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtk_mock"
	"github.com/coyim/gotk3adapter/gtki"
	mck "github.com/stretchr/testify/mock"
)

type mockedGTK struct {
	mm mck.Mock
	gtk_mock.Mock
}

func (m *mockedGTK) SettingsGetDefault() (gtki.Settings, error) {
	args := m.mm.Called()

	var ret gtki.Settings
	retv := args.Get(0)
	if retv != nil {
		ret = retv.(gtki.Settings)
	}

	return ret, args.Error(1)
}

func (m *mockedGTK) BuilderNew() (gtki.Builder, error) {
	args := m.mm.Called()

	var ret gtki.Builder
	retv := args.Get(0)
	if retv != nil {
		ret = retv.(gtki.Builder)
	}

	return ret, args.Error(1)
}

type mockedSettings struct {
	mm mck.Mock
	gtk_mock.MockSettings
}

func (m *mockedSettings) GetProperty(v string) (interface{}, error) {
	args := m.mm.Called(v)

	return args.Get(0), args.Error(1)
}

type mockedBuilder struct {
	mm mck.Mock
	gtk_mock.MockBuilder
}

func (m *mockedBuilder) GetObject(v string) (glibi.Object, error) {
	args := m.mm.Called(v)

	var ret glibi.Object
	retv := args.Get(0)
	if retv != nil {
		ret = retv.(glibi.Object)
	}

	return ret, args.Error(1)
}

type mockedListBox struct {
	mm mck.Mock
	gtk_mock.MockListBox
}

func (m *mockedListBox) GetStyleContext() (gtki.StyleContext, error) {
	args := m.mm.Called()

	var ret gtki.StyleContext
	retv := args.Get(0)
	if retv != nil {
		ret = retv.(gtki.StyleContext)
	}

	return ret, args.Error(1)
}

type mockedStyleContext struct {
	mm mck.Mock
	gtk_mock.MockStyleContext
}

func (m *mockedStyleContext) GetProperty2(v string, v2 gtki.StateFlags) (interface{}, error) {
	args := m.mm.Called(v, v2)

	return args.Get(0), args.Error(1)
}
