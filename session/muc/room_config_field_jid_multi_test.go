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
		c.Assert(field.Raw(), DeepEquals, mock.raw)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}

func (*MucRoomConfigFieldJidMultiSuite) Test_RoomConfigFieldJidMultiValue_SetValue(c *C) {
	field := newRoomConfigFieldJidMultiValue([]string{"bla"})
	c.Assert(field.Raw(), Equals, []jid.Any{jid.Parse("bla")})
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.SetValue(nil)
	c.Assert(field.Raw(), Equals, []jid.Any{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValue([]string{"bla@domain.org"})
	c.Assert(field.Raw(), Equals, []jid.Any{jid.Parse("bla@domain.org")})
	c.Assert(field.Value(), DeepEquals, []string{"bla@domain.org"})

	field.SetValue([]string{"foo@domain.org", "invalid*@whatever"})
	c.Assert(field.Raw(), Equals, []jid.Any{jid.Parse("foo@domain.org")})
	c.Assert(field.Value(), DeepEquals, []string{"foo@domain.org"})
}
