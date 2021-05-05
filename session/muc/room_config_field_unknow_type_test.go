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
		c.Assert(field.Raw(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}
