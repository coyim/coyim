package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/gotk3/gotk3/gtk"
)

type application struct {
	*gliba.Application
	internal *gtk.Application
}

func WrapApplicationSimple(v *gtk.Application) gtki.Application {
	if v == nil {
		return nil
	}
	return &application{gliba.WrapApplicationSimple(&v.Application), v}
}

func WrapApplication(v *gtk.Application, e error) (gtki.Application, error) {
	return WrapApplicationSimple(v), e
}

func UnwrapApplication(v gtki.Application) *gtk.Application {
	if v == nil {
		return nil
	}
	return v.(*application).internal
}

func (v *application) GetActiveWindow() gtki.Window {
	ret := WrapWindowSimple(v.internal.GetActiveWindow())
	if ret == nil {
		return nil
	}
	return ret
}
