package gui

import (
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigInfo struct {
	content gtki.Grid `gtk-widget:"room-config-info-content"`
}

func (rc *roomConfigAssistant) newRoomConfigInfo() *roomConfigInfo {
	ri := &roomConfigInfo{}

	ri.initBuilder()
	return ri
}

func (ri *roomConfigInfo) initBuilder() {
	b := newBuilder("MUCRoomConfigInfo")
	panicOnDevError(b.bindObjects(ri))
}
