package gui

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtk_mock"
	"github.com/coyim/gotk3adapter/gtki"
	. "gopkg.in/check.v1"
)

type FirstAccountSuite struct{}

var _ = Suite(&FirstAccountSuite{})

type mockBuilderNew struct {
	gtk_mock.Mock
	builderToReturn gtki.Builder
}

type mockBuilderWithGetObjectAndConnectSignals struct {
	mockBuilder

	getObjectReturn glibi.Object
	getObjectArg    string

	connectSignalsArgs map[string]interface{}
}

func (m *mockBuilderWithGetObjectAndConnectSignals) GetObject(v string) (glibi.Object, error) {
	m.getObjectArg = v
	return m.getObjectReturn, nil
}

func (m *mockBuilderWithGetObjectAndConnectSignals) ConnectSignals(v map[string]interface{}) {
	m.connectSignalsArgs = v
}

func (v *mockBuilderNew) BuilderNew() (gtki.Builder, error) {
	return v.builderToReturn, nil
}

type mockDialogWithShowAll struct {
	gtk_mock.MockDialog

	showAllCalled bool
}

func (m *mockDialogWithShowAll) ShowAll() {
	m.showAllCalled = true
}

func (*FirstAccountSuite) Test_thatDialogGetsShown(c *C) {
	orgG := g
	defer func() {
		g = orgG
	}()
	g = Graphics{}

	m := &glibIdleAddMock{}
	g.glib = m

	idleAddWaitChannel := make(chan bool)

	m.f = func(ff interface{}, vals ...interface{}) (glibi.SourceHandle, error) {
		go func() {
			ff.(func())()
			idleAddWaitChannel <- true
		}()
		return glibi.SourceHandle(0), nil
	}

	gtkm := &mockBuilderNew{}
	mockDialog := &mockDialogWithShowAll{}
	mb := &mockBuilderWithGetObjectAndConnectSignals{
		getObjectReturn: mockDialog,
	}
	gtkm.builderToReturn = mb
	g.gtk = gtkm

	ui := &gtkUI{}

	methodDone := make(chan bool)

	// Real test functionality starts here

	go func() {
		ui.showFirstAccountWindow()
		methodDone <- true
	}()

	<-idleAddWaitChannel
	c.Assert(mockDialog.showAllCalled, Equals, true)
	c.Assert(mb.getObjectArg, Equals, "dialog")
	c.Assert(mb.stringGiven, Matches, "(?s).*?Setup your first account.*?")

	mb.connectSignalsArgs["on_cancel"].(func())()

	<-idleAddWaitChannel
	<-idleAddWaitChannel

	<-methodDone
}
