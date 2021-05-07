package muc

import xmppData "github.com/coyim/coyim/xmpp/data"

// RoomConfigFieldOptions contains information of the option of
// of the field part of the configuration form
type RoomConfigFieldOption struct {
	Value string
	Label string
}

func newRoomConfigFieldOption(v, l string) *RoomConfigFieldOption {
	return &RoomConfigFieldOption{v, l}
}

func formFieldOptionsValues(options []xmppData.FormFieldOptionX) []*RoomConfigFieldOption {
	l := []*RoomConfigFieldOption{}
	for _, o := range options {
		l = append(l, newRoomConfigFieldOption(o.Value, o.Label))
	}
	return l
}
