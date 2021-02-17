package gui

import (
	"time"

	"github.com/coyim/gotk3adapter/gtki"
)

type infoBarComponent struct {
	text            string
	messageType     gtki.MessageType
	canBeClosed     bool
	onCloseCallback func()

	infoBar    gtki.InfoBar `gtk-widget:"infobar"`
	timeBox    gtki.Box     `gtk-widget:"time-box"`
	timeLabel  gtki.Label   `gtk-widget:"time-label"`
	titleLabel gtki.Label   `gtk-widget:"title-label"`
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
			select {
			case <-ticker.C:
				doInUIThread(func() {
					ib.timeLabel.SetText(elapsedFriendlyTime(t))
				})
			}
		}
	}()
	ib.timeLabel.SetText(elapsedFriendlyTime(t))
	ib.timeBox.Show()
}
