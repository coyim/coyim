package gdk_mock

import "github.com/coyim/gotk3adapter/gdki"

func init() {
	gdki.AssertGdk(&Mock{})
	gdki.AssertEvent(&MockEvent{})
	gdki.AssertEventButton(&MockEventButton{})
	gdki.AssertEventKey(&MockEventKey{})
	gdki.AssertPixbuf(&MockPixbuf{})
	gdki.AssertPixbufLoader(&MockPixbufLoader{})
	gdki.AssertScreen(&MockScreen{})
	gdki.AssertWindow(&MockWindow{})
}
