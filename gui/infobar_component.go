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
	text                   string
	messageType            gtki.MessageType
	doWhenRequestedToClose func() // doWhenRequestedToClose will be called from the UI thread
	tickerCancelChannel    chan bool

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
				if ib.doWhenRequestedToClose != nil {
					ib.doWhenRequestedToClose()
				}

				go ib.closeActiveTickerChannel()
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

// whenRequestedToClose MUST be called from the UI thread
func (ib *infoBarComponent) whenRequestedToClose(doWhenClose func()) {
	ib.doWhenRequestedToClose = doWhenClose
	ib.showCloseButton(ib.doWhenRequestedToClose != nil)
}

// showCloseButton MUST be called from the UI thread
func (ib *infoBarComponent) showCloseButton(v bool) {
	ib.infoBar.SetShowCloseButton(v)
}

// addActionWidget MUST be called from the UI thread
func (ib *infoBarComponent) addActionWidget(w gtki.Widget, responseType gtki.ResponseType) {
	ib.infoBar.AddActionWidget(w, responseType)
	ib.infoBar.ShowAll()
}

func (ib *infoBarComponent) view() gtki.InfoBar {
	return ib.infoBar
}

const infoBarTimeFormat = "January 2, 2006 at 15:04:05"

func (ib *infoBarComponent) setTime(t time.Time) {
	ib.refreshElapsedTime(t)
	ib.time.Show()

	friendlyTime := formatTimeWithLayout(t, infoBarTimeFormat)
	ib.time.SetTooltipText(friendlyTime)

	go ib.tickNotificationTime(t)
}

// tickNotificationTime MUST NOT be called from the UI thread
func (ib *infoBarComponent) tickNotificationTime(t time.Time) {
	ib.tickerCancelChannel = make(chan bool)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			doInUIThread(func() {
				ib.refreshElapsedTime(t)
			})
		case <-ib.tickerCancelChannel:
			return
		}
	}
}

// closeActiveTickerChannel MUST NOT be called from the UI thread
func (ib *infoBarComponent) closeActiveTickerChannel() {
	if ib.tickerCancelChannel != nil {
		close(ib.tickerCancelChannel)
	}
}

// refreshElapsedTime MUST be called from the UI thread
func (ib *infoBarComponent) refreshElapsedTime(t time.Time) {
	ib.time.SetText(elapsedFriendlyTime(t))
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
