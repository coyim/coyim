package gui

import (
	"time"

	"github.com/coyim/gotk3adapter/gtki"
)

type infoBarType int

const (
	infoBarTypeInfo infoBarType = iota
	infoBarTypeWarning
	infoBarTypeQuestion
	infoBarTypeError
	infoBarTypeOther
)

const (
	infoBarInfoIconName     = "message_info"
	infoBarWarningIconName  = "message_warning"
	infoBarQuestionIconName = "message_question"
	infoBarErrorIconName    = "message_error"
)

var infoBarIconNames = map[infoBarType]string{
	infoBarTypeInfo:     infoBarInfoIconName,
	infoBarTypeWarning:  infoBarWarningIconName,
	infoBarTypeQuestion: infoBarQuestionIconName,
	infoBarTypeError:    infoBarErrorIconName,
}

type infoBarComponent struct {
	text            string
	messageType     gtki.MessageType
	canBeClosed     bool
	onCloseCallback func()

	infoBar    gtki.InfoBar `gtk-widget:"infobar"`
	timeBox    gtki.Box     `gtk-widget:"time-box"`
	iconTime   gtki.Image   `gtk-widget:"icon-time"`
	timeLabel  gtki.Label   `gtk-widget:"time-label"`
	titleLabel gtki.Label   `gtk-widget:"title-label"`
	icon       gtki.Image   `gtk-widget:"icon-image"`
}

func (u *gtkUI) newInfoBarComponent(text string, messageType gtki.MessageType) *infoBarComponent {
	ib := &infoBarComponent{
		text:        text,
		messageType: messageType,
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

	ib.titleLabel.SetText(ib.text)
	ib.infoBar.SetMessageType(ib.messageType)
	mucStyles.setInfoBarStyle(ib.infoBar)

	tp := infoBarTypeForMessageType(messageType)
	if icoName, ok := infoBarIconNames[tp]; ok {
		ib.icon.SetFromPixbuf(getMUCIconPixbuf(icoName))
		ib.icon.Show()
	}

	return ib
}

// setClosable MUST be called from the UI thread
func (ib *infoBarComponent) setClosable(v bool) {
	ib.canBeClosed = v
	ib.infoBar.SetShowCloseButton(v)
}

// addActionWidget MUST be called from the UI thread
func (ib *infoBarComponent) addActionWidget(w gtki.Widget, responseType gtki.ResponseType) {
	ib.infoBar.AddActionWidget(w, responseType)
	ib.infoBar.ShowAll()
}

func (ib *infoBarComponent) isClosable() bool {
	return ib.canBeClosed
}

func (ib *infoBarComponent) onClose(f func()) {
	ib.onCloseCallback = f
}

func (ib *infoBarComponent) view() gtki.InfoBar {
	return ib.infoBar
}

func (ib *infoBarComponent) setTickerTime(t time.Time) {
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for {
			<-ticker.C
			doInUIThread(func() {
				ib.timeLabel.SetText(elapsedFriendlyTime(t))
			})
		}
	}()

	ib.timeLabel.SetText(elapsedFriendlyTime(t))
	ib.timeBox.Show()
}

func infoBarTypeForMessageType(mt gtki.MessageType) infoBarType {
	switch mt {
	case gtki.MESSAGE_INFO:
		return infoBarTypeInfo
	case gtki.MESSAGE_WARNING:
		return infoBarTypeWarning
	case gtki.MESSAGE_QUESTION:
		return infoBarTypeQuestion
	case gtki.MESSAGE_ERROR:
		return infoBarTypeError
	}
	return infoBarTypeOther
}
