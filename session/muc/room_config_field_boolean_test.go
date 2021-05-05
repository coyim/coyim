package muc

import (
	"strconv"

	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldBooleanSuite struct{}

var _ = Suite(&MucRoomConfigFieldBooleanSuite{})

func (*MucRoomConfigFieldBooleanSuite) Test_newRoomConfigFieldBooleanValue(c *C) {
	cases := []struct {
		values   []string
		expected bool
	}{
		{
			[]string{"true"},
			true,
		},
		{
			[]string{"1"},
			true,
		},
		{
			[]string{"true", "bla", "foo"},
			true,
		},
		{
			[]string{"True"},
			true,
		},
		{
			[]string{"false"},
			false,
		},
		{
			[]string{"0"},
			false,
		},
		{
			[]string{"False"},
			false,
		},
		{
			[]string{"whatever"},
			false,
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldBooleanValue(mock.values)
		c.Assert(field.Raw(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, []string{strconv.FormatBool(mock.expected)})
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldBooleanValue_SetValue(c *C) {
	field := newRoomConfigFieldBooleanValue([]string{"false"})
	c.Assert(field.Raw(), Equals, false)
	c.Assert(field.Value(), DeepEquals, []string{"false"})

	field.SetValue(true)
	c.Assert(field.Raw(), Equals, true)
	c.Assert(field.Value(), DeepEquals, []string{"true"})

	field.SetValue("true")
	c.Assert(field.Raw(), Equals, true)
	c.Assert(field.Value(), DeepEquals, []string{"true"})

	field.SetValue("bla")
	c.Assert(field.Raw(), Equals, true)
	c.Assert(field.Value(), DeepEquals, []string{"true"})

	field.SetValue(false)
	c.Assert(field.Raw(), Equals, false)
	c.Assert(field.Value(), DeepEquals, []string{"false"})
}
