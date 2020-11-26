package gui

import "github.com/coyim/gotk3adapter/gtki"

type mucRoomConfigComponent struct {
	notebook gtki.Notebook `gtk-widget:"config-room-notebook"`
}

func (u *gtkUI) newMUCRoomConfigComponent() *mucRoomConfigComponent {
	c := &mucRoomConfigComponent{}

	b := newBuilder("MUCRoomConfig")
	panicOnDevError(b.bindObjects(c))

	return c
}

func (c *mucRoomConfigComponent) configurationView() gtki.Widget {
	return c.notebook
}
