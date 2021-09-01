package gui

import (
	"sync"
	"time"

	"github.com/coyim/coyim/i18n"
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
	infoBarOtherIconName    = "message_other"
)

var infoBarIconNames = map[infoBarType]string{
	infoBarTypeInfo:     infoBarInfoIconName,
	infoBarTypeWarning:  infoBarWarningIconName,
	infoBarTypeQuestion: infoBarQuestionIconName,
	infoBarTypeError:    infoBarErrorIconName,
	infoBarTypeOther:    infoBarOtherIconName,
}

type infoBarComponent struct {
	u                      *gtkUI
	text                   string
	messageType            gtki.MessageType
	spinner                *spinner
	doWhenRequestedToClose func() // doWhenRequestedToClose will be called from the UI thread
	tickerCancelChannel    chan bool
	tickerCancelLock       sync.Mutex

	infoBar    gtki.InfoBar `gtk-widget:"infobar"`
	time       gtki.Label   `gtk-widget:"time-label"`
	title      gtki.Label   `gtk-widget:"title-label"`
	icon       gtki.Image   `gtk-widget:"icon-image"`
	spinnerBox gtki.Box     `gtk-widget:"spinner-box"`
}

func (u *gtkUI) newInfoBarComponent(text string, messageType gtki.MessageType) *infoBarComponent {
	ib := &infoBarComponent{
		u:           u,
		text:        text,
		messageType: messageType,
	}

	ib.initBuilder()
	ib.initSpinnerComponent()
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

func (ib *infoBarComponent) initSpinnerComponent() {
	ib.spinner = ib.u.newSpinnerComponent()
	ib.spinner.setSize(spinnerSize24)
	ib.spinnerBox.Add(ib.spinner.spinner())
}

func (ib *infoBarComponent) initDefaults() {
	formatter := newInfobarHighlightFormatter(ib.text)
	formatter.formatLabel(ib.title)

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

// addAction MUST be called from the UI thread
func (ib *infoBarComponent) addAction(label string, responseType gtki.ResponseType, signals map[string]interface{}) {
	action, _ := g.gtk.ButtonNewWithLabel(label)
	for signal, handler := range signals {
		action.Connect(signal, handler)
	}

	ib.addActionWidget(action, responseType)
}

// addActionWidget MUST be called from the UI thread
func (ib *infoBarComponent) addActionWidget(w gtki.Widget, responseType gtki.ResponseType) {
	ib.infoBar.AddActionWidget(w, responseType)
	w.Show()
}

// reveal MUST be called from the UI thread
func (ib *infoBarComponent) reveal() {
	ib.infoBar.ShowAll()
	g.gtk.InfoBarSetRevealed(ib.infoBar, true)
}

func (ib *infoBarComponent) view() gtki.InfoBar {
	return ib.infoBar
}

// setTime MUST be called from the UI thread
func (ib *infoBarComponent) setTime(t time.Time) {
	ib.refreshElapsedTime(t)
	ib.time.Show()

	friendlyTime := formatTimeWithLayout(t, i18n.Local("January 2, 2006 at 15:04:05"))
	ib.time.SetTooltipText(friendlyTime)

	go ib.tickNotificationTime(t)
}

// tickNotificationTime MUST NOT be called from the UI thread
func (ib *infoBarComponent) tickNotificationTime(t time.Time) {
	ib.tickerCancelLock.Lock()
	ib.tickerCancelChannel = make(chan bool)
	ib.tickerCancelLock.Unlock()

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
	ib.tickerCancelLock.Lock()
	defer ib.tickerCancelLock.Unlock()

	if ib.tickerCancelChannel != nil {
		close(ib.tickerCancelChannel)
	}
}

// refreshElapsedTime MUST be called from the UI thread
func (ib *infoBarComponent) refreshElapsedTime(t time.Time) {
	ib.time.SetText(elapsedFriendlyTime(t))
}

// showSpinner MUST be called from the UI thread
func (ib *infoBarComponent) showSpinner() {
	ib.spinner.show()
	ib.spinnerBox.Show()
}

// hideSpinner MUST be called from the UI thread
func (ib *infoBarComponent) hideSpinner() {
	ib.spinner.hide()
	ib.spinnerBox.Hide()
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
