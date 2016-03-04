package gui

import (
	"github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/gotk3adapter/glibi"
	. "gopkg.in/check.v1"
)

type SignalsSuite struct{}

var _ = Suite(&SignalsSuite{})

type mockSignal struct {
	glib_mock.MockSignal
	name string
}

type mockGlibSignalNew struct {
	glib_mock.Mock
}

func (*mockGlibSignalNew) SignalNew(v1 string) (glibi.Signal, error) {
	return &mockSignal{name: v1}, nil
}

func (*SignalsSuite) Test_initSignals_willInitTheSignals(c *C) {
	g = Graphics{glib: &mockGlibSignalNew{}}
	initSignals()
	c.Assert(accountChangedSignal.(*mockSignal).name, Equals, "coyim-account-changed")
	c.Assert(enableWindow.(*mockSignal).name, Equals, "enable")
	c.Assert(disableWindow.(*mockSignal).name, Equals, "disable")
}
