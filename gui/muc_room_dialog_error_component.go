package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type dialogErrorOptions struct {
	title        string
	header       string
	message      string
	parentWindow gtki.Window
}

// dialogErrorComponent is shown after an error occurred.
type dialogErrorComponent struct {
	title   string
	header  string
	message string

	dialog       gtki.Dialog `gtk-widget:"room-error-dialog"`
	headerLabel  gtki.Label  `gtk-widget:"room-error-dialog-header"`
	messageLabel gtki.Label  `gtk-widget:"room-error-dialog-message"`

	parentWindow gtki.Window
}

func createDialogErrorComponent(options dialogErrorOptions) *dialogErrorComponent {
	d := &dialogErrorComponent{
		title:        options.title,
		header:       options.header,
		message:      options.message,
		parentWindow: options.parentWindow,
	}

	d.initBuilder()
	d.initDefaults()

	return d
}

func (d *dialogErrorComponent) initBuilder() {
	builder := newBuilder("MUCRoomDialogErrorComponent")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_ok": d.onOkClicked,
	})
}

func (d *dialogErrorComponent) initDefaults() {
	mucStyles.setRoomDialogErrorComponentHeaderStyle(d.headerLabel)

	d.dialog.SetTitle(d.title)
	d.headerLabel.SetText(d.header)
	d.messageLabel.SetText(d.message)

	if d.parentWindow != nil {
		d.dialog.SetTransientFor(d.parentWindow)
	}
}

// onOkClicked MUST be called from the UI thread
func (d *dialogErrorComponent) onOkClicked() {
	d.dialog.Destroy()
}

// setParent MUST be called from the UI thread
func (d *dialogErrorComponent) setParent(p gtki.Window) {
	d.dialog.SetTransientFor(p)
}

// show MUST be called from the UI thread
func (d *dialogErrorComponent) show() {
	d.dialog.Show()
}

// updateMessageError MUST be called from the UI thread
func (d *dialogErrorComponent) updateMessageError(msg string) {
	d.messageLabel.SetText(msg)
}
