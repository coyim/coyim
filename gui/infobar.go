package gui

import "github.com/coyim/gotk3adapter/gtki"

type infoBar struct {
	text            string
	messageType     gtki.MessageType
	isClosable      bool
	onCloseCallback func()

	widget  gtki.InfoBar `gtk-widget:"bar"`
	content gtki.Box     `gtk-widget:"content"`
	title   gtki.Label   `gtk-widget:"title"`
}

func newInfoBar(text string, mt gtki.MessageType) *infoBar {
	ib := &infoBar{
		text:            text,
		messageType:     mt,
		onCloseCallback: func() {},
	}

	builder := newBuilder("InfoBar")
	panicOnDevError(builder.bindObjects(ib))

	builder.ConnectSignals(map[string]interface{}{
		"on_close": func() {
			if ib.isClosable {
				ib.onCloseCallback()
			}
		},
	})

	ib.title.SetText(ib.text)
	ib.widget.SetMessageType(ib.messageType)

	return ib
}

func (ib *infoBar) setClosable(v bool) {
	ib.isClosable = v
	ib.widget.SetShowCloseButton(v)
}

func (ib *infoBar) onClose(f func()) {
	ib.onCloseCallback = f
}

func (ib *infoBar) addAction(text string, rt gtki.ResponseType) {
	ib.widget.AddButton(text, rt)
}

// getMessageType implements the "infoMessage" interface
func (ib *infoBar) getMessageType() gtki.MessageType {
	return ib.messageType
}

// getWidget implements the "widget" interface
func (ib *infoBar) getWidget() gtki.Widget {
	return ib.widget
}
