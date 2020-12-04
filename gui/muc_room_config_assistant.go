package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigAssistant struct {
	assistant gtki.Assistant `gtk-widget:"room-config-assistant"`
}

func newRoomConfigAssistant() *roomConfigAssistant {
	rc := &roomConfigAssistant{}
	rc.initBuilder()
	return rc
}

func (rc *roomConfigAssistant) initBuilder() {
	b := newBuilder("MUCRoomConfigAssistant")
	panicOnDevError(b.bindObjects(rc))

	b.ConnectSignals(map[string]interface{}{
		"on_cancel":       rc.onCancel,
		"on_page_changed": rc.onPageChanged,
	})
}

func (rc *roomConfigAssistant) onCancel() {
	rc.assistant.Destroy()
}

func (rc *roomConfigAssistant) onPageChanged(_ gtki.Assistant, pg gtki.Widget) {
	rc.assistant.SetPageComplete(pg, true)
	switch rc.assistant.GetCurrentPage() {
	case 0:
		// TODO: Add implementation for "room information" step
	case 1:
		// TODO: Add implementation for "room access" step
	case 2:
		// TODO: Add implementation for "room permissions" step
	case 3:
		// TODO: Add implementation for "room occupants" step
	case 4:
		// TODO: Add implementation for "room others configurations" step
	case 5:
		// TODO: Add implementation for "summary configurations" step
	}
}

func (rc *roomConfigAssistant) show() {
	rc.assistant.ShowAll()
}
