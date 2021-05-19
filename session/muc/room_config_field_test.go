package muc

import (
	xmppData "github.com/coyim/coyim/xmpp/data"
	"github.com/coyim/coyim/xmpp/jid"
	. "gopkg.in/check.v1"
)

type MucRoomConfigFieldsSuite struct{}

var _ = Suite(&MucRoomConfigFieldsSuite{})

func (*MucRoomConfigFieldsSuite) Test_RoomConfigFormField_roomConfigFormFieldValueFactory(c *C) {
	field := &RoomConfigFormField{}

	cases := []struct {
		fieldX           xmppData.FormFieldX
		expectedValue    []string
		expectedRawValue interface{}
	}{
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldText,
				Values: []string{"bla"},
			},
			[]string{"bla"},
			"bla",
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldTextMulti,
				Values: []string{"bla", "foo"},
			},
			[]string{"bla", "foo"},
			[]string{"bla", "foo"},
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldBoolean,
				Values: []string{"true"},
			},
			[]string{"true"},
			true,
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldList,
				Values: []string{"foo"},
			},
			[]string{"foo"},
			"foo",
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldListMulti,
				Values: []string{"one", "two"},
			},
			[]string{"one", "two"},
			[]string{"one", "two"},
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldJidMulti,
				Values: []string{"foo", "bla@domain.org"},
			},
			[]string{"foo", "bla@domain.org"},
			[]jid.Any{jid.Parse("foo"), jid.Parse("bla@domain.org")},
		},
		{
			xmppData.FormFieldX{
				Type:   "unknow type",
				Values: []string{"foo", "bla"},
			},
			[]string{"foo", "bla"},
			[]string{"foo", "bla"},
		},
	}

	for _, fieldCase := range cases {
		field.value = roomConfigFormFieldValueFactory(fieldCase.fieldX)
		c.Assert(field.Value(), DeepEquals, fieldCase.expectedValue)
		c.Assert(field.RawValue(), DeepEquals, fieldCase.expectedRawValue)
	}
}

func (*MucRoomConfigFieldsSuite) Test_RoomConfigFormField_ValueType(c *C) {
	field := &RoomConfigFormField{}

	cases := []struct {
		fieldX    xmppData.FormFieldX
		valueType HasRoomConfigFormFieldValue
	}{
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldText,
				Values: []string{"foo", "bla"},
			},
			newRoomConfigFieldTextValue([]string{"foo"}),
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldTextMulti,
				Values: []string{"one", "two", "three"},
			},
			newRoomConfigFieldTextMultiValue([]string{"one", "two", "three"}),
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldBoolean,
				Values: []string{"foo", "bla"},
			},
			newRoomConfigFieldBooleanValue([]string{"foo", "bla"}),
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldList,
				Values: []string{"foo", "bla"},
			},
			newRoomConfigFieldListValue([]string{"foo", "bla"}, formFieldOptionsValues(nil)),
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldListMulti,
				Values: []string{"foo", "bla"},
			},
			newRoomConfigFieldListMultiValue([]string{"foo", "bla"}, formFieldOptionsValues(nil)),
		},
		{
			xmppData.FormFieldX{
				Type:   RoomConfigFieldJidMulti,
				Values: []string{"foo", "bla"},
			},
			newRoomConfigFieldJidMultiValue([]string{"foo", "bla"}),
		},
		{
			xmppData.FormFieldX{
				Type:   "unknow field type",
				Values: []string{"foo", "bla"},
			},
			newRoomConfigFieldUnknownValue([]string{"foo", "bla"}),
		},
	}

	for _, fieldCase := range cases {
		field.value = roomConfigFormFieldValueFactory(fieldCase.fieldX)
		c.Assert(field.ValueType(), DeepEquals, fieldCase.valueType)
	}
}
