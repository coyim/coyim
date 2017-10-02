package gdka

import "github.com/coyim/gotk3adapter/gdki"

func init() {
	gdki.AssertGdk(&RealGdk{})
	gdki.AssertEvent(&event{})
	gdki.AssertEventButton(&eventButton{})
	gdki.AssertEventKey(&eventKey{})
	gdki.AssertPixbuf(&pixbuf{})
	gdki.AssertPixbufLoader(&pixbufLoader{})
	gdki.AssertScreen(&screen{})
	gdki.AssertWindow(&window{})
}
