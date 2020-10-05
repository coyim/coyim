package gui

import "github.com/coyim/gotk3adapter/gtki"

type roomViewLoadingInfoBar struct {
	parent gtki.Box

	infoBar gtki.InfoBar `gtk-widget:"info"`
}

func (v *roomView) newRoomViewLoadingInfoBar(parent gtki.Box) *roomViewLoadingInfoBar {
	ib := &roomViewLoadingInfoBar{
		parent: parent,
	}

	builder := newBuilder("MUCRoomLoadingInfoBar")
	panicOnDevError(builder.bindObjects(ib))

	ib.parent.Add(ib.infoBar)

	return ib
}

func (ib *roomViewLoadingInfoBar) show() {
	ib.infoBar.Show()
}

func (ib *roomViewLoadingInfoBar) hide() {
	ib.infoBar.Hide()
}

func (ib *roomViewLoadingInfoBar) getWidget() gtki.InfoBar {
	return ib.infoBar
}
