package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldTextSuite struct{}

var _ = Suite(&MucRoomConfigFieldTextSuite{})

func (*MucRoomConfigFieldTextSuite) Test_newRoomConfigFieldTextValue(c *C) {
	cases := []struct {
		values   []string
		expected string
	}{
		{
			[]string{"bla"},
			"bla",
		},
		{
			[]string{"foo"},
			"foo",
		},
		{
			[]string{""},
			"",
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldTextValue(mock.values)
		c.Assert(field.Raw(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, []string{mock.expected})
	}
}

func (*MucRoomConfigFieldTextSuite) Test_RoomConfigFieldTextValue_SetValue(c *C) {
	field := newRoomConfigFieldTextValue([]string{"false"})
	c.Assert(field.Raw(), Equals, "false")
	c.Assert(field.Value(), DeepEquals, []string{"false"})

	field.SetValue("bla")
	c.Assert(field.Raw(), Equals, "bla")
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.SetValue("foo")
	c.Assert(field.Raw(), Equals, "foo")
	c.Assert(field.Value(), DeepEquals, []string{"foo"})

	field.SetValue(nil)
	c.Assert(field.Raw(), Equals, "foo")
	c.Assert(field.Value(), DeepEquals, []string{"foo"})

	field.SetValue("")
	c.Assert(field.Raw(), Equals, "")
	c.Assert(field.Value(), DeepEquals, []string{""})
}
