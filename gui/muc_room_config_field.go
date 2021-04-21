package gui

import (
	"errors"

	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type hasRoomConfigFormField interface {
	fieldWidget() gtki.Widget
	fieldName() string
	fieldLabel() string
	fieldValue() interface{}
}

var (
	errRoomConfigFieldNotSupported = errors.New("room configuration form field not supported")
)

func roomConfigFormFieldFactory(field *muc.RoomConfigFormField) (hasRoomConfigFormField, error) {
	switch field.Type {
	case muc.RoomConfigFieldText:
		return newRoomConfigFormTextField(field), nil
	}

	return nil, errRoomConfigFieldNotSupported
}
