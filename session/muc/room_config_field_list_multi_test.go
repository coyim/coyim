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
		field := newRoomConfigFieldListMultiValue(mock.values, nil)
		c.Assert(field.Selected(), DeepEquals, mock.expected)
		c.Assert(field.Value(), DeepEquals, mock.expected)
	}
}

func (*MucRoomConfigFieldListMultiSuite) Test_RoomConfigFieldListMultiValue_SetValue(c *C) {
	field := newRoomConfigFieldListMultiValue([]string{"bla", "Juan", "dog"}, nil)
	c.Assert(field.Selected(), DeepEquals, []string{"bla", "Juan", "dog"})
	c.Assert(field.Value(), DeepEquals, []string{"bla", "Juan", "dog"})

	field.SetSelected([]string{})
	c.Assert(field.Selected(), DeepEquals, []string{})
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetSelected(nil)
	c.Assert(field.Selected(), IsNil)
	c.Assert(field.Value(), IsNil)

	field.SetSelected([]string{"foooooo"})
	c.Assert(field.Selected(), DeepEquals, []string{"foooooo"})
	c.Assert(field.Value(), DeepEquals, []string{"foooooo"})
}

func (*MucRoomConfigFieldListMultiSuite) Test_RoomConfigFieldListMultiValue_optionsWorks(c *C) {
	field := newRoomConfigFieldListMultiValue(nil, []*RoomConfigFieldOption{{"one", "One"}, {"two", "Two"}})
	c.Assert(field.Options(), DeepEquals, []*RoomConfigFieldOption{{"one", "One"}, {"two", "Two"}})

	field.SetOptions([]*RoomConfigFieldOption{{"bla", "Bla"}, {"foo", "Foo"}})
	c.Assert(field.Options(), DeepEquals, []*RoomConfigFieldOption{{"bla", "Bla"}, {"foo", "Foo"}})

	field.SetOptions(nil)
	c.Assert(field.Options(), DeepEquals, []*RoomConfigFieldOption{{"bla", "Bla"}, {"foo", "Foo"}})
}

func (*MucRoomConfigFieldListMultiSuite) Test_RoomConfigFieldListMultiValue_IsSelected(c *C) {
	field := newRoomConfigFieldListMultiValue([]string{"one"}, []*RoomConfigFieldOption{{"one", "One"}, {"two", "Two"}, {"three", "Three"}})
	c.Assert(field.IsSelected("one"), Equals, true)
	c.Assert(field.IsSelected("two"), Equals, false)
	c.Assert(field.IsSelected("something"), Equals, false)

	field.SetSelected([]string{"three"})
	c.Assert(field.IsSelected("three"), Equals, true)
	c.Assert(field.IsSelected("one"), Equals, false)
	c.Assert(field.IsSelected("two"), Equals, false)
	c.Assert(field.IsSelected("something"), Equals, false)

	field.SetSelected([]string{"one", "three"})
	c.Assert(field.IsSelected("three"), Equals, true)
	c.Assert(field.IsSelected("one"), Equals, true)
	c.Assert(field.IsSelected("two"), Equals, false)
	c.Assert(field.IsSelected("something"), Equals, false)

	field.SetSelected(nil)
	c.Assert(field.IsSelected("three"), Equals, false)
	c.Assert(field.IsSelected("one"), Equals, false)
	c.Assert(field.IsSelected("two"), Equals, false)
	c.Assert(field.IsSelected("something"), Equals, false)
}
