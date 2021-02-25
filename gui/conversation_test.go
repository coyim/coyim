package gui

import (
	"github.com/coyim/coyim/i18n"
	"github.com/coyim/gotk3adapter/glib_mock"

	. "gopkg.in/check.v1"
)

type GUIXmppSuite struct{}

var _ = Suite(&GUIXmppSuite{})

func (s *GUIXmppSuite) Test_createStatusMessage_createsStatusMessages(c *C) {
	i18n.InitLocalization(&glib_mock.Mock{})

	c.Assert(createStatusMessage("Foo", "", "", false), Equals, "Foo is now Available")

	c.Assert(createStatusMessage("Foo", "", "", true), Equals, "Foo is now Offline")
	c.Assert(createStatusMessage("Foo", "", "something happened", true), Equals, "Foo is now Offline (something happened)")
	c.Assert(createStatusMessage("Foo", "xa", "something happened", true), Equals, "Foo is now Offline (Not Available: something happened)")
	c.Assert(createStatusMessage("Foo", "xa", "", true), Equals, "Foo is now Offline (Not Available)")

	c.Assert(createStatusMessage("Foo", "xa", "something happened", false), Equals, "Foo is now Not Available (something happened)")

	c.Assert(createStatusMessage("Foo2", "dnd", "", false), Equals, "Foo2 is now Busy")
	c.Assert(createStatusMessage("Foo2", "dnd", "some stuff", false), Equals, "Foo2 is now Busy (some stuff)")
	c.Assert(createStatusMessage("Foo2", "some state", "", false), Equals, "Foo2 is now some state")

	c.Assert(createStatusMessage("Foo2", "chat", "", false), Equals, "Foo2 is now Free for Chat")
	c.Assert(createStatusMessage("Foo2", "chat", "really!", false), Equals, "Foo2 is now Free for Chat (really!)")

	c.Assert(createStatusMessage("Foo3", "away", "", false), Equals, "Foo3 is now Away")
	c.Assert(createStatusMessage("Foo3", "away", "wut", false), Equals, "Foo3 is now Away (wut)")

	c.Assert(createStatusMessage("Foo4", "invisible", "", false), Equals, "Foo4 is now Invisible")
	c.Assert(createStatusMessage("Foo4", "invisible", "wut", false), Equals, "Foo4 is now Invisible (wut)")
}
