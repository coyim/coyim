package muc

import xmppData "github.com/coyim/coyim/xmpp/data"

// RoomConfigFieldOption contains information of the field option, part of the configuration form
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
		l = append(l, newRoomConfigFieldOption(o.Value, o.Var))
	}
	return l
}
