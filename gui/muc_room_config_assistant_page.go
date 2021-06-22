package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomConfigAssistantPage struct {
	page    gtki.ScrolledWindow `gtk-widget:"room-config-assistant-page"`
	content gtki.Box            `gtk-widget:"room-config-assistant-page-content"`
}

func newRoomConfigAssistantPage(p *roomConfigPage) *roomConfigAssistantPage {
	ap := &roomConfigAssistantPage{}

	builder := newBuilder("MUCRoomConfigAssistantePage")
	panicOnDevError(builder.bindObjects(ap))

	ap.content.Add(p.page)

	return ap
}
