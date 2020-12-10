package gui

import (
	"github.com/coyim/coyim/coylog"
	"github.com/coyim/coyim/session/muc"
	"github.com/coyim/gotk3adapter/gtki"
)

type roomConfigInfo struct {
	form *muc.RoomConfigForm

	content        gtki.Grid        `gtk-widget:"room-config-info-content"`
	title          gtki.Label       `gtk-widget:"room-title"`
	instructions   gtki.Label       `gtk-widget:"room-instructions"`
	roomName       gtki.Entry       `gtk-widget:"roomname"`
	roomDesc       gtki.TextView    `gtk-widget:"roomdesc"`
	roomLang       gtki.Entry       `gtk-widget:"lang"`
	roomPersistent gtki.CheckButton `gtk-widget:"persistentroom"`
	roomPublicroom gtki.CheckButton `gtk-widget:"publicroom"`

	log coylog.Logger
}

func (rc *roomConfigAssistant) newRoomConfigInfo(form *muc.RoomConfigForm) *roomConfigInfo {
	ri := &roomConfigInfo{
		form: form,
	}

	ri.initBuilder()
	return ri
}

func (ri *roomConfigInfo) initBuilder() {
	b := newBuilder("MUCRoomConfigInfo")
	panicOnDevError(b.bindObjects(ri))
}
