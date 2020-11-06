package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomDestroyView struct {
	spinner *spinner

	transient  gtki.Window
	dialog     gtki.Dialog `gtk-widget:"destroy-room-dialog"`
	spinnerBox gtki.Box    `gtk-widget:"destroy-room-spinner-box"`

	cancelChannel chan bool
}

func newRoomDestroyView(t gtki.Window) *roomDestroyView {
	d := &roomDestroyView{
		transient: t,
	}

	d.initBuilder()
	d.initDefaults()

	return d
}

func (d *roomDestroyView) initBuilder() {
	builder := newBuilder("MUCRoomDestroyDialog")
	panicOnDevError(builder.bindObjects(d))

	builder.ConnectSignals(map[string]interface{}{
		"on_destroy_clicked":  d.onDestroyRoom,
		"on_cancel_clicked":   d.onCancel,
		"on_dialog_destroyed": d.onDialogDestroy,
	})
}

func (d *roomDestroyView) initDefaults() {
	d.dialog.SetTransientFor(d.transient)

	d.spinner = newSpinner()
	d.spinnerBox.Add(d.spinner.getWidget())
}

func (d *roomDestroyView) onDestroyRoom() {
	// TODO: Implement the logic behind this functionality and how
	// the view interact with the user
	d.spinner.show()
}

func (d *roomDestroyView) onCancel() {
	d.spinner.hide()
	d.cancelAnyRequest()
	d.close()
}

func (d *roomDestroyView) onDialogDestroy() {
	d.cancelAnyRequest()
}

func (d *roomDestroyView) cancelAnyRequest() {
	if d.cancelChannel != nil {
		d.cancelChannel <- true
	}
}

func (d *roomDestroyView) show() {
	d.dialog.Show()
}

func (d *roomDestroyView) close() {
	d.dialog.Destroy()
}
