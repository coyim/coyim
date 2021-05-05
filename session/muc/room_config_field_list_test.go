package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldListSuite struct{}

var _ = Suite(&MucRoomConfigFieldListSuite{})

func (*MucRoomConfigFieldListSuite) Test_newRoomConfigFieldListValue(c *C) {
	cases := []struct {
		values   []string
		expected string
	}{
		{
			[]string{},
			"",
		},
		{
			[]string{"whatever"},
			"whatever",
		},
		{
			[]string{"bla", "foo"},
			"bla",
		},
		{
			[]string{"foo", "bla"},
			"foo",
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldListValue(mock.values)
		c.Assert(field.Raw(), Equals, mock.expected)
		c.Assert(field.Value(), DeepEquals, []string{mock.expected})
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldListValue_SetValue(c *C) {
	field := newRoomConfigFieldListValue([]string{"bla", "foo"})
	c.Assert(field.Raw(), DeepEquals, "bla")
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.SetValue("foo")
	c.Assert(field.Raw(), DeepEquals, "foo")
	c.Assert(field.Value(), DeepEquals, []string{"foo"})

	field.SetValue("whatever")
	c.Assert(field.Raw(), DeepEquals, "whatever")
	c.Assert(field.Value(), DeepEquals, []string{"whatever"})

	field.SetValue("")
	c.Assert(field.Raw(), DeepEquals, "")
	c.Assert(field.Value(), DeepEquals, []string{""})

	field.SetValue(1000)
	c.Assert(field.Raw(), DeepEquals, "")
	c.Assert(field.Value(), DeepEquals, []string{""})

	field.SetValue("abc")
	c.Assert(field.Raw(), DeepEquals, "abc")
	c.Assert(field.Value(), DeepEquals, []string{"abc"})
}
