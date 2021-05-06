package muc

import (
	xmppData "github.com/coyim/coyim/xmpp/data"
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

func (*MucRoomConfigFieldsSuite) Test_RoomConfigFormField(c *C) {
	field := &RoomConfigFormField{}

	field.value = roomConfigFormFieldValueFactory(xmppData.FormFieldX{
		Type:   RoomConfigFieldBoolean,
		Values: []string{"true"},
	})
	c.Assert(field.RawValue(), DeepEquals, true)
	c.Assert(field.Value(), DeepEquals, []string{"true"})

	field.value = roomConfigFormFieldValueFactory(xmppData.FormFieldX{
		Type:   RoomConfigFieldText,
		Values: []string{"bla"},
	})
	c.Assert(field.RawValue(), Equals, "bla")
	c.Assert(field.Value(), DeepEquals, []string{"bla"})

	field.value = roomConfigFormFieldValueFactory(xmppData.FormFieldX{
		Type:   RoomConfigFieldList,
		Values: []string{"foo"},
	})
	c.Assert(field.RawValue(), DeepEquals, "foo")
	c.Assert(field.Value(), DeepEquals, []string{"foo"})
}
