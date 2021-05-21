package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldTextMultiSuite struct{}

var _ = Suite(&MucRoomConfigFieldBooleanSuite{})

func (*MucRoomConfigFieldBooleanSuite) Test_newRoomConfigFieldTextMultiValue(c *C) {
	cases := []struct {
		values        []string
		expectedValue []string
		expectedText  string
	}{
		{
			[]string{},
			[]string{},
			"",
		},
		{
			[]string{"bla", "foo"},
			[]string{"bla", "foo"},
			"bla foo",
		},
		{
			[]string{"my", "multiline", "text"},
			[]string{"my", "multiline", "text"},
			"my multiline text",
		},
	}

	for _, mock := range cases {
		field := newRoomConfigFieldTextMultiValue(mock.values)
		c.Assert(field.Text(), DeepEquals, mock.expectedText)
		c.Assert(field.Value(), DeepEquals, mock.expectedValue)
	}
}

func (*MucRoomConfigFieldBooleanSuite) Test_RoomConfigFieldTextMultiValue_SetValue(c *C) {
	field := newRoomConfigFieldTextMultiValue([]string{"bla", "foo"})
	c.Assert(field.Text(), DeepEquals, "bla foo")
	c.Assert(field.Value(), DeepEquals, []string{"bla", "foo"})

	field.SetText([]string{"text"})
	c.Assert(field.Text(), DeepEquals, "text")
	c.Assert(field.Value(), DeepEquals, []string{"text"})

	field.SetText([]string{"whatever"})
	c.Assert(field.Text(), DeepEquals, "whatever")
	c.Assert(field.Value(), DeepEquals, []string{"whatever"})

	field.SetText([]string{})
	c.Assert(field.Text(), DeepEquals, "")
	c.Assert(field.Value(), DeepEquals, []string{})

	field.SetText(nil)
	c.Assert(field.Text(), Equals, "")
	c.Assert(field.Value(), IsNil)
}
