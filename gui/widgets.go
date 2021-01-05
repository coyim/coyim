package gui

import "github.com/coyim/gotk3adapter/gtki"

type widget interface {
	widget() gtki.Widget
}

type message interface {
	messageType() gtki.MessageType
}
