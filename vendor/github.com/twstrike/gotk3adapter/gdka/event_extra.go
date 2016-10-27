package gdka

import "github.com/twstrike/gotk3adapter/gdki"

func eventCast(e gdki.Event) *event {
	if e == nil {
		return nil
	}
	return e.(*event)
}
