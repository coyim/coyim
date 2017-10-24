package gtka

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
)

type application struct {
	*gliba.Application
	internal *gtk.Application
}

func wrapApplicationSimple(v *gtk.Application) *application {
	if v == nil {
		return nil
	}
	return &application{gliba.WrapApplicationSimple(&v.Application), v}
}

func wrapApplication(v *gtk.Application, e error) (*application, error) {
	return wrapApplicationSimple(v), e
}

func unwrapApplication(v gtki.Application) *gtk.Application {
	if v == nil {
		return nil
	}
	return v.(*application).internal
}

func (v *application) GetActiveWindow() gtki.Window {
	ret := wrapWindowSimple(v.internal.GetActiveWindow())
	if ret == nil {
		return nil
	}
	return ret
}
