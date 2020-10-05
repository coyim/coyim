package gui

import (
	"sync"

	"github.com/coyim/gotk3adapter/gtki"
)

type roomViewWarning struct {
	text string
	lock sync.Mutex

	bar     gtki.Box   `gtk-widget:"warning-infobar"`
	message gtki.Label `gtk-widget:"message"`
}

func newRoomViewWarning(text string) *roomViewWarning {
	w := &roomViewWarning{
		text: text,
	}

	builder := newBuilder("MUCRoomWarning")
	panicOnDevError(builder.bindObjects(w))

	w.message.SetText(w.text)

	return w
}

type roomViewWarningsInfoBar struct {
	infoBar gtki.InfoBar `gtk-widget:"bar"`
}

func (v *roomView) newRoomViewWarningsInfoBar(onShowWarnings func(), onClose func()) *roomViewWarningsInfoBar {
	ib := &roomViewWarningsInfoBar{}

	builder := newBuilder("MUCRoomWarningsInfoBar")
	panicOnDevError(builder.bindObjects(ib))

	builder.ConnectSignals(map[string]interface{}{
		"on_show_warnings": onShowWarnings,
		"on_close":         onClose,
	})

	return ib
}

// getMessageType implements the "message" interface
func (ib *roomViewWarningsInfoBar) getMessageType() gtki.MessageType {
	return gtki.MESSAGE_WARNING
}

// getWidget implements the "widget" interface
func (ib *roomViewWarningsInfoBar) getWidget() gtki.Widget {
	return ib.infoBar
}

type roomViewWarningsOverlay struct {
	warnings []*roomViewWarning
	onClose  func()

	box      gtki.Box      `gtk-widget:"warningsBox"`
	revealer gtki.Revealer `gtk-widget:"revealer"`
}

func (v *roomView) newRoomViewWarningsOverlay(onClose func()) *roomViewWarningsOverlay {
	o := &roomViewWarningsOverlay{
		onClose: func() {
			if onClose != nil {
				onClose()
			}
		},
	}

	builder := newBuilder("MUCRoomWarningsOverlay")
	panicOnDevError(builder.bindObjects(o))

	builder.ConnectSignals(map[string]interface{}{
		"on_close": o.close,
	})

	prov := providerWithStyle("box", style{
		"padding": "12px",
	})

	updateWithStyle(o.box, prov)

	v.messagesBox.Add(o.revealer)

	return o
}

func (o *roomViewWarningsOverlay) add(text string) {
	w := newRoomViewWarning(text)
	o.warnings = append(o.warnings, w)

	prov := providerWithStyle("box", style{
		"color":            "#744210",
		"background-color": "#fefcbf",
		"border":           "1px solid #d69e2e",
		"border-radius":    "4px",
		"padding":          "10px",
	})

	updateWithStyle(w.bar, prov)

	o.box.PackStart(w.bar, false, false, 5)

	o.box.ShowAll()
}

func (o *roomViewWarningsOverlay) show() {
	o.revealer.SetRevealChild(true)
}

func (o *roomViewWarningsOverlay) hide() {
	o.revealer.SetRevealChild(false)
}

func (o *roomViewWarningsOverlay) close() {
	o.hide()
	o.onClose()
}

func (o *roomViewWarningsOverlay) clear() {
	// TODO: Why can't we just remove
	// all the entitites inside the warningsArea
	// and then remove the need to have the "warnings" field at all
	for _, w := range o.warnings {
		o.box.Remove(w.bar)
	}
}
