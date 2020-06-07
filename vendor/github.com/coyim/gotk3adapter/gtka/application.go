package gtka

import (
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/coyim/gotk3adapter/glibi"
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

func (v *application) AddWindow(w gtki.Window) {
	v.internal.AddWindow(UnwrapWindow(w))
}

func (v *application) RemoveWindow(w gtki.Window) {
	v.internal.RemoveWindow(UnwrapWindow(w))
}

func (v *application) PrefersAppMenu() bool {
	return v.internal.PrefersAppMenu()
}

func (v *application) GetAppMenu() glibi.MenuModel {
	return gliba.WrapMenuModelSimple(v.internal.GetAppMenu())
}

func (v *application) SetAppMenu(val glibi.MenuModel) {
	v.internal.SetAppMenu(gliba.UnwrapMenuModel(val))
}

func (v *application) GetMenubar() glibi.MenuModel {
	return gliba.WrapMenuModelSimple(v.internal.GetMenubar())
}

func (v *application) SetMenubar(val glibi.MenuModel) {
	v.internal.SetMenubar(gliba.UnwrapMenuModel(val))
}
