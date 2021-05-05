package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldListMultiSuite struct{}

var _ = Suite(&MucRoomConfigFieldListMultiSuite{})

func (*MucRoomConfigFieldListMultiSuite) Test_newRoomConfigFieldListMultiValue(c *C) {
	cases := []struct {
		values   []string
		expected []string
	}{
		{
			[]string{"bla", "foo"},
			[]string{"bla", "foo"},
		},
		{
			[]string{"1", "2", "3"},
			[]string{"1", "2", "3"},
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldListMultiValue(mock.values)
		c.Assert(field.Raw(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}

func (*MucRoomConfigFieldListMultiSuite) Test_RoomConfigFieldListMultiValue_SetValue(c *C) {
	field := newRoomConfigFieldListMultiValue([]string{"bla", "Juan", "dog"})
	c.Assert(field.Raw(), DeepEquals, []string{"bla", "Juan", "dog"})
	c.Assert(field.Value(), DeepEquals, []string{"bla", "Juan", "dog"})

	field.SetValue([]string{})
	c.Assert(field.Raw(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValue(nil)
	c.Assert(field.Raw(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValue("")
	c.Assert(field.Raw(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValue([]string{"foooooo"})
	c.Assert(field.Raw(), DeepEquals, []string{"foooooo"})
	c.Assert(field.Value(), DeepEquals, []string{"foooooo"})
}
