package muc

import (
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldJidMultiSuite struct{}

var _ = Suite(&MucRoomConfigFieldBooleanSuite{})

func (*MucRoomConfigFieldJidMultiSuite) Test_newRoomConfigFieldJidMultiValue(c *C) {
	cases := []struct {
		values   []string
		raw      []jid.Any
		expected []string
	}{
		{
			[]string{"bla", "foo"},
			[]jid.Any{jid.Parse("bla"), jid.Parse("foo")},
			[]string{"bla", "foo"},
		},
		{
			[]string{"bla", "@@@@"},
			[]jid.Any{jid.Parse("bla")},
			[]string{"bla"},
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldJidMultiValue(mock.values)
		c.Assert(field.List(), DeepEquals, mock.raw)
		c.Assert(field.Value(), DeepEquals, mock.expected)
		c.Assert(field.Length(), Equals, len(mock.expected))
	}
}

func (*MucRoomConfigFieldJidMultiSuite) Test_RoomConfigFieldJidMultiValue_setsRightValue(c *C) {
	field := newRoomConfigFieldJidMultiValue([]string{"bla"})
	c.Assert(field.List(), Equals, []jid.Any{jid.Parse("bla")})
	c.Assert(field.Value(), DeepEquals, []string{"bla"})
	c.Assert(field.Length(), Equals, 1)

	field.SetValue(nil)
	c.Assert(field.List(), Equals, []jid.Any{})
	c.Assert(field.Value(), DeepEquals, []string{})
	c.Assert(field.Length(), Equals, 0)

	field.SetValue([]string{"bla@domain.org"})
	c.Assert(field.List(), Equals, []jid.Any{jid.Parse("bla@domain.org")})
	c.Assert(field.Value(), DeepEquals, []string{"bla@domain.org"})
	c.Assert(field.Length(), Equals, 1)

	field.SetValue([]string{"foo@domain.org", "invalid*@whatever"})
	c.Assert(field.List(), Equals, []jid.Any{jid.Parse("foo@domain.org")})
	c.Assert(field.Value(), DeepEquals, []string{"foo@domain.org"})
	c.Assert(field.Length(), Equals, 1)
}

func (*MucRoomConfigFieldJidMultiSuite) Test_RoomConfigFieldJidMultiValue_initValues(c *C) {
	field := newRoomConfigFieldJidMultiValue([]string{"bla", "foo"})
	c.Assert(field.List(), DeepEquals, []jid.Any{jid.Parse("bla"), jid.Parse("foo")})
	c.Assert(field.Length(), Equals, 2)

	field.initValues([]string{"whatever", "bla@domain.org"})
	c.Assert(field.List(), DeepEquals, []jid.Any{jid.Parse("whatever"), jid.Parse("bla@domain.org")})
	c.Assert(field.Length(), Equals, 2)

	field.initValues(nil)
	c.Assert(field.List(), IsNil)
	c.Assert(field.Length(), Equals, 0)
}

func (*MucRoomConfigFieldJidMultiSuite) Test_RoomConfigFieldJidMultiValue_returnTheRightLength(c *C) {
	field := newRoomConfigFieldJidMultiValue([]string{"bla", "foo"})
	c.Assert(field.Length(), Equals, 2)

	field.SetValue([]string{"bla"})
	c.Assert(field.Length(), Equals, 1)

	field.SetValue(nil)
	c.Assert(field.Length(), Equals, 0)

	field.SetValue([]string{"foo"})
	c.Assert(field.Length(), Equals, 1)

	field.SetValue([]string{})
	c.Assert(field.Length(), Equals, 0)
}
