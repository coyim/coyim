package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldUnknowTypeSuite struct{}

var _ = Suite(&MucRoomConfigFieldUnknowTypeSuite{})

func (*MucRoomConfigFieldUnknowTypeSuite) Test_newRoomConfigFieldUnknowValue(c *C) {
	cases := []struct {
		values   []string
		expected []string
	}{
		{
			[]string{},
			[]string{},
		},
		{
			[]string{"bla", "foo"},
			[]string{"bla", "foo"},
		},
		{
			[]string{"whatever", ""},
			[]string{"whatever", ""},
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldUnknowValue(mock.values)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}

func (*MucRoomConfigFieldUnknowTypeSuite) Test_RoomConfigFieldUnknowValue_ValueAndSetValue(c *C) {
	field := newRoomConfigFieldUnknowValue([]string{"foo", "bla"})
	c.Assert(field.Value(), DeepEquals, []string{"foo", "bla"})

	field.SetValue([]string{"foo"})
	c.Assert(field.Value(), DeepEquals, []string{"foo"})

	field.SetValue(nil)
	c.Assert(field.Value(), IsNil)

	field.SetValue([]string{"1", "2", "3"})
	c.Assert(field.Value(), DeepEquals, []string{"1", "2", "3"})
}
