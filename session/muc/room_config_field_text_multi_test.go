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
		c.Assert(field.Text(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldTextMultiValue_SetValue(c *C) {
	field := newRoomConfigFieldTextMultiValue([]string{"bla", "foo"})
	c.Assert(field.Text(), DeepEquals, []string{"bla", "foo"})
	c.Assert(field.Value(), DeepEquals, []string{"bla", "foo"})

	field.SetText([]string{"text"})
	c.Assert(field.Text(), DeepEquals, []string{"text"})
	c.Assert(field.Value(), DeepEquals, []string{"text"})

	field.SetText([]string{"whatever"})
	c.Assert(field.Text(), DeepEquals, []string{"whatever"})
	c.Assert(field.Value(), DeepEquals, []string{"whatever"})

	field.SetText([]string{})
	c.Assert(field.Text(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetText(nil)
	c.Assert(field.Text(), IsNil)
	c.Assert(field.Value(), DeepEquals, []string{})
}
