package gui

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glib_mock"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/glibi"
	. "github.com/twstrike/coyim/Godeps/_workspace/src/gopkg.in/check.v1"
)

type UIThreadSuite struct{}

var _ = Suite(&UIThreadSuite{})

type glibIdleAddMock struct {
	glib_mock.Mock
	f func(interface{}, ...interface{}) (glibi.SourceHandle, error)
}

func (v *glibIdleAddMock) IdleAdd(v1 interface{}, v2 ...interface{}) (glibi.SourceHandle, error) {
	return v.f(v1, v2...)
}

type glibMainDepthMock struct {
	glib_mock.Mock
	mainDepth int
}

func (v *glibMainDepthMock) MainDepth() int {
	return v.mainDepth
}

func (*UIThreadSuite) Test_assertInUIThread_panicsIfNotInUIThread(c *C) {
	m := &glibMainDepthMock{mainDepth: 1}
	g = Graphics{glib: m}

	assertInUIThread()

	m.mainDepth = 0
	c.Assert(func() {
		assertInUIThread()
	}, Panics, "This function has to be called from the UI thread")
}

func (*UIThreadSuite) Test_doInUIThread(c *C) {
	m := &glibIdleAddMock{}
	g = Graphics{glib: m}

	m.f = func(ff interface{}, vals ...interface{}) (glibi.SourceHandle, error) {
		ffx := ff.(func())
		ffx()
		return glibi.SourceHandle(0), nil
	}

	ran := false
	doInUIThread(func() {
		ran = true
	})
	c.Assert(ran, Equals, true)
}
