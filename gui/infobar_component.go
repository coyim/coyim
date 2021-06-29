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
	onCloseCallback func() // onCloseCallback will be called from the UI thread

	infoBar gtki.InfoBar `gtk-widget:"infobar"`
	time    gtki.Label   `gtk-widget:"time-label"`
	title   gtki.Label   `gtk-widget:"title-label"`
	icon    gtki.Image   `gtk-widget:"icon-image"`
}

func (u *gtkUI) newInfoBarComponent(text string, messageType gtki.MessageType) *infoBarComponent {
	ib := &infoBarComponent{
		text:        text,
		messageType: messageType,
	}

	ib.initBuilder()
	ib.initDefaults()
	ib.initStyleAndIcon()

	return ib
}

func (ib *infoBarComponent) initBuilder() {
	builder := newBuilder("InfoBar")
	panicOnDevError(builder.bindObjects(ib))

	builder.ConnectSignals(map[string]interface{}{
		"handle-response": func(info gtki.InfoBar, response gtki.ResponseType) {
			if response == gtki.RESPONSE_CLOSE {
				if ib.canBeClosed && ib.onCloseCallback != nil {
					ib.onCloseCallback()
				}
			}
		},
	})
}

func (ib *infoBarComponent) initDefaults() {
	ib.title.SetText(ib.text)
	ib.infoBar.SetMessageType(ib.messageType)
}

func (ib *infoBarComponent) initStyleAndIcon() {
	mucStyles.setInfoBarStyle(ib.infoBar)

	tp := infoBarTypeForMessageType(ib.messageType)
	if icoName, ok := infoBarIconNames[tp]; ok {
		ib.icon.SetFromPixbuf(getMUCIconPixbuf(icoName))
		ib.icon.Show()
	}

	if actions, err := ib.infoBar.GetActionArea(); err == nil {
		actions.SetProperty("margin", 0)
	}
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
	ib.setClosable(f != nil)
}

func (ib *infoBarComponent) view() gtki.InfoBar {
	return ib.infoBar
}

const infoBarTimeFormat = "January 2, 2006 at 15:04:05"

func (ib *infoBarComponent) setTime(t time.Time) {
	friendlyTime := formatTimeWithLayout(t, infoBarTimeFormat)
	ib.time.SetTooltipText(friendlyTime)

	ib.refreshElapsedTime(t)

	go ib.tickNotificationTime(t)
}

// tickNotificationTime MUST NOT be called from the UI thread
func (ib *infoBarComponent) tickNotificationTime(t time.Time) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			doInUIThread(func() {
				ib.refreshElapsedTime(t)
			})
		}
	}
}

// refreshElapsedTime MUST be called from the UI thread
func (ib *infoBarComponent) refreshElapsedTime(t time.Time) {
	ib.time.SetText(elapsedFriendlyTime(t))
	ib.time.Show()
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
