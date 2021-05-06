package muc

import "github.com/coyim/coyim/xmpp/jid"

// RoomConfigFieldJidMultiValue contains information of the value of the jid multi field
type RoomConfigFieldJidMultiValue struct {
	value []jid.Any
}

func newRoomConfigFieldJidMultiValue(values []string) *RoomConfigFieldJidMultiValue {
	v := &RoomConfigFieldJidMultiValue{}

	v.SetValues(values)

	return v
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldJidMultiValue) Value() []string {
	value := []string{}
	for _, addr := range v.value {
		value = append(value, addr.String())
	}
	return value
}

// SetValues will try to set the new list based on the provided values
func (v *RoomConfigFieldJidMultiValue) SetValues(values []string) {
	v.value = []jid.Any{}
	for _, addr := range values {
		if any := jid.Parse(addr); any.Valid() {
			v.value = append(v.value, any)
		}
	}
}

// List return the current list of jids
func (v *RoomConfigFieldJidMultiValue) List() []jid.Any {
	return v.value
}
