package gui

import "github.com/coyim/gotk3adapter/gtki"

type withWidget interface {
	widget() gtki.Widget
}

type withMessage interface {
	messageType() gtki.MessageType
}
