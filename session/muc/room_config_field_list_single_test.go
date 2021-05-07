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
		field := newRoomConfigFieldListValue(mock.values, nil)
		c.Assert(field.Selected(), Equals, mock.expected)
		c.Assert(field.Value(), DeepEquals, []string{mock.expected})
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldListValue_SetValue(c *C) {
	field := newRoomConfigFieldListValue([]string{"bla", "foo"}, nil)
	c.Assert(field.Selected(), DeepEquals, "bla")
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.SetSelected("foo")
	c.Assert(field.Selected(), DeepEquals, "foo")
	c.Assert(field.Value(), DeepEquals, []string{"foo"})

	field.SetSelected("whatever")
	c.Assert(field.Selected(), DeepEquals, "whatever")
	c.Assert(field.Value(), DeepEquals, []string{"whatever"})

	field.SetSelected("")
	c.Assert(field.Selected(), DeepEquals, "")
	c.Assert(field.Value(), DeepEquals, []string{""})

	field.SetSelected("abc")
	c.Assert(field.Selected(), DeepEquals, "abc")
	c.Assert(field.Value(), DeepEquals, []string{"abc"})
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldListValue_Options(c *C) {
	field := newRoomConfigFieldListValue(nil, []*RoomConfigFieldOption{{"bla", "foo"}})
	c.Assert(field.Options(), DeepEquals, []*RoomConfigFieldOption{{"bla", "foo"}})

	field.SetOptions(nil)
	c.Assert(field.Options(), IsNil)

	field.SetOptions([]*RoomConfigFieldOption{{Value: "foo"}})
	c.Assert(field.Options(), DeepEquals, []*RoomConfigFieldOption{{Value: "foo"}})

	field.SetOptions([]*RoomConfigFieldOption{{Value: "whatever"}})
	c.Assert(field.Options(), DeepEquals, []*RoomConfigFieldOption{{Value: "whatever"}})
}
