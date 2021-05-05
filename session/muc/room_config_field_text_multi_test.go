package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldTextMultiSuite struct{}

var _ = Suite(&MucRoomConfigFieldBooleanSuite{})

func (*MucRoomConfigFieldBooleanSuite) Test_newRoomConfigFieldTextMultiValue(c *C) {
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
			[]string{"my", "multiline", "text"},
			[]string{"my", "multiline", "text"},
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldTextMultiValue(mock.values)
		c.Assert(field.Raw(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldTextMultiValue_SetValue(c *C) {
	field := newRoomConfigFieldTextMultiValue([]string{"bla", "foo"})
	c.Assert(field.Raw(), DeepEquals, []string{"bla", "foo"})
	c.Assert(field.Value(), DeepEquals, []string{"bla", "foo"})

	field.SetValue([]string{"text"})
	c.Assert(field.Raw(), DeepEquals, []string{"text"})
	c.Assert(field.Value(), DeepEquals, []string{"text"})

	field.SetValue([]string{"whatever"})
	c.Assert(field.Raw(), DeepEquals, []string{"whatever"})
	c.Assert(field.Value(), DeepEquals, []string{"whatever"})

	field.SetValue([]string{})
	c.Assert(field.Raw(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValue("bla")
	c.Assert(field.Raw(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetValue([]string{"multi", "line", "text"})
	c.Assert(field.Raw(), DeepEquals, []string{"multi", "line", "text"})
	c.Assert(field.Value(), DeepEquals, []string{"multi", "line", "text"})

	field.SetValue(20000)
	c.Assert(field.Raw(), DeepEquals, []string{"multi", "line", "text"})
	c.Assert(field.Value(), DeepEquals, []string{"multi", "line", "text"})
}
