package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type infoBar struct {
	text            string
	mt              gtki.MessageType
	canBeClosed     bool
	onCloseCallback func()

	bar     gtki.InfoBar `gtk-widget:"bar"`
	content gtki.Box     `gtk-widget:"content"`
	title   gtki.Label   `gtk-widget:"title"`
}

func (u *gtkUI) newInfoBarComponent(text string, mt gtki.MessageType) *infoBar {
	ib := &infoBar{
		text: text,
		mt:   mt,
	}

	builder := newBuilder("InfoBar")
	panicOnDevError(builder.bindObjects(ib))

	builder.ConnectSignals(map[string]interface{}{
		"handle-response": func(info gtki.InfoBar, response gtki.ResponseType) {
			if response != gtki.RESPONSE_CLOSE {
				return
			}

			if ib.canBeClosed && ib.onCloseCallback != nil {
				ib.onCloseCallback()
			}
		},
	})

	ib.title.SetText(ib.text)
	ib.bar.SetMessageType(ib.mt)

	return ib
}

// setClosable MUST be called from the UI thread
func (ib *infoBar) setClosable(v bool) {
	ib.canBeClosed = v
	ib.bar.SetShowCloseButton(v)
}

func (ib *infoBar) isClosable() bool {
	return ib.canBeClosed
}

func (ib *infoBar) onClose(f func()) {
	ib.onCloseCallback = f
}

// messageType implements the "withMessage" interface
func (ib *infoBar) messageType() gtki.MessageType {
	return ib.mt
}

// widget implements the "withWidget" interface
func (ib *infoBar) widget() gtki.Widget {
	return ib.bar
}
