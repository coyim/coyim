package muc

import (
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldsSuite struct{}

var _ = Suite(&MucRoomConfigFieldsSuite{})

func (*MucRoomConfigFieldsSuite) Test_newConfigListSingleField(c *C) {
	val := []string{"a", "b"}
	res := newConfigListSingleField(val)
	c.Assert(res.(*configListSingleField).options, DeepEquals, val)
}

func (*MucRoomConfigFieldsSuite) Test_configListSingleField_UpdateField(c *C) {
	vv := &configListSingleField{}
	vv.UpdateField("hello", nil)
	c.Assert(vv.value, Equals, "hello")

	vv.UpdateField("goodbye", []string{"something", "else"})
	c.Assert(vv.CurrentValue(), Equals, "goodbye")
	c.Assert(vv.Options(), DeepEquals, []string{"something", "else"})
}

func (*MucRoomConfigFieldsSuite) Test_newConfigListMultiField(c *C) {
	val := []string{"a", "b"}
	res := newConfigListMultiField(val)
	c.Assert(res.(*configListMultiField).options, DeepEquals, val)
}

func (*MucRoomConfigFieldsSuite) Test_configListMultiField_UpdateField(c *C) {
	vv := &configListMultiField{}
	vv.UpdateField([]string{"hello"}, []string{"something", "else"})
	c.Assert(vv.values, DeepEquals, []string{"hello"})
	c.Assert(vv.options, DeepEquals, []string{"something", "else"})
}

func (*MucRoomConfigFieldsSuite) Test_RoomConfigFormField(c *C) {
	field := &RoomConfigFormField{}

	field.SetValue(true)
	field.Type = RoomConfigFieldBoolean
	c.Assert(field.Value, Equals, true)
	c.Assert(field.GetValue(), DeepEquals, []string{"true"})

	field.SetValue("bla")
	c.Assert(field.Value, Equals, "bla")

	field.Type = RoomConfigFieldText
	c.Assert(field.GetValue(), DeepEquals, []string{"bla"})

	field.SetValue([]string{"foo"})
	c.Assert(field.Value, DeepEquals, []string{"foo"})
}
