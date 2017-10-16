package gliba

import (
	"github.com/coyim/gotk3adapter/glibi"
	"github.com/gotk3/gotk3/glib"
)

type Application struct {
	*Object
	*glib.Application
}

func WrapApplicationSimple(v *glib.Application) *Application {
	if v == nil {
		return nil
	}
	return &Application{WrapObjectSimple(v.Object), v}
}

func unwrapApplication(v glibi.Application) *glib.Application {
	if v == nil {
		return nil
	}
	return v.(*Application).Application
}
