package gtk_mock

import "github.com/coyim/gotk3adapter/gtki"

type MockInfoBar struct {
	MockBox
}

func (*MockInfoBar) GetOrientation() gtki.Orientation {
	return gtki.HorizontalOrientation
}

func (*MockInfoBar) SetOrientation(o gtki.Orientation) {
}

func (*MockInfoBar) AddActionWidget(w gtki.Widget, responseId gtki.ResponseType) {
}

func (*MockInfoBar) AddButton(buttonText string, responseId gtki.ResponseType) {
}

func (*MockInfoBar) SetResponseSensitive(responseId gtki.ResponseType, setting bool) {
}

func (*MockInfoBar) SetDefaultResponse(responseId gtki.ResponseType) {
}

func (*MockInfoBar) SetMessageType(messageType gtki.MessageType) {
}

func (*MockInfoBar) GetMessageType() gtki.MessageType {
	return gtki.MESSAGE_OTHER
}

func (*MockInfoBar) GetActionArea() (gtki.Widget, error) {
	return nil, nil
}

func (*MockInfoBar) GetContentArea() (gtki.Box, error) {
	return nil, nil
}

func (*MockInfoBar) GetShowCloseButton() bool {
	return false
}

func (*MockInfoBar) SetShowCloseButton(setting bool) {
}
