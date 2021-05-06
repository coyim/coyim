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
	}
}

func (*MucRoomConfigFieldJidMultiSuite) Test_RoomConfigFieldJidMultiValue_setsRightValue(c *C) {
	field := newRoomConfigFieldJidMultiValue([]string{"bla"})
	c.Assert(field.List(), Equals, []jid.Any{jid.Parse("bla")})
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.SetValues(nil)
	c.Assert(field.List(), Equals, []jid.Any{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValues([]string{"bla@domain.org"})
	c.Assert(field.List(), Equals, []jid.Any{jid.Parse("bla@domain.org")})
	c.Assert(field.Value(), DeepEquals, []string{"bla@domain.org"})

	field.SetValues([]string{"foo@domain.org", "invalid*@whatever"})
	c.Assert(field.List(), Equals, []jid.Any{jid.Parse("foo@domain.org")})
	c.Assert(field.Value(), DeepEquals, []string{"foo@domain.org"})
}

func (*MucRoomConfigFieldJidMultiSuite) Test_RoomConfigFieldJidMultiValue_initValues(c *C) {
	field := newRoomConfigFieldJidMultiValue([]string{"bla", "foo"})
	c.Assert(field.List(), DeepEquals, []jid.Any{jid.Parse("bla"), jid.Parse("foo")})

	field.SetValues([]string{"whatever", "bla@domain.org"})
	c.Assert(field.List(), DeepEquals, []jid.Any{jid.Parse("whatever"), jid.Parse("bla@domain.org")})

	field.SetValues(nil)
	c.Assert(field.List(), IsNil)
}
