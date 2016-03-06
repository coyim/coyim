package gui

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

type UIThreadSuite struct{}

var _ = Suite(&UIThreadSuite{})

func (*UIThreadSuite) Test_assertInUIThread_panicsIfNotInUIThread(c *C) {
	weAreInUIThread = true
	assertInUIThread()
	weAreInUIThread = false
	c.Assert(func() {
		assertInUIThread()
	}, Panics, "This function have to be called from the UI thread")
}

type glibIdleAddMock struct {
	glib_mock.Mock
	f func(interface{}, ...interface{}) (glibi.SourceHandle, error)
}

func (v *glibIdleAddMock) IdleAdd(v1 interface{}, v2 ...interface{}) (glibi.SourceHandle, error) {
	return v.f(v1, v2...)
}

func (*UIThreadSuite) Test_doInUIThread(c *C) {
	m := &glibIdleAddMock{}
	g = Graphics{glib: m}

	m.f = func(ff interface{}, vals ...interface{}) (glibi.SourceHandle, error) {
		ffx := ff.(func())
		ffx()
		return glibi.SourceHandle(0), nil
	}

	c.Assert(weAreInUIThread, Equals, false)
	ran := false
	doInUIThread(func() {
		ran = true
		c.Assert(weAreInUIThread, Equals, true)
	})
	c.Assert(ran, Equals, true)
	c.Assert(weAreInUIThread, Equals, false)
}
