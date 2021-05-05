package muc

import "github.com/coyim/coyim/xmpp/jid"

// RoomConfigFieldJidMultiValue contains information of the value of the jid multi field
type RoomConfigFieldJidMultiValue struct {
	value []jid.Any
}

func newRoomConfigFieldJidMultiValue(values []string) HasRoomConfigFormFieldValue {
	v := &RoomConfigFieldJidMultiValue{}

	v.initValues(values)

	return v
}

func (v *RoomConfigFieldJidMultiValue) initValues(values []string) {
	for _, addr := range values {
		if any := jid.Parse(addr); any.Valid() {
			v.value = append(v.value, any)
		}
	}
}

// Value implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldJidMultiValue) Value() []string {
	value := []string{}
	for _, addr := range v.value {
		value = append(value, addr.String())
	}
	return value
}

// SetValue implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldJidMultiValue) SetValue(value interface{}) {
	if val, ok := value.([]string); ok {
		v.initValues(val)
	}
}

// Raw implements the "HasRoomConfigFormFieldValue" interface
func (v *RoomConfigFieldJidMultiValue) Raw() interface{} {
	return v.value
}
