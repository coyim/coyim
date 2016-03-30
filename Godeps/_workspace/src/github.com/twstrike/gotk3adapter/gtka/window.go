package gtka

import (
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/gotk3/gotk3/gtk"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdka"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gdki"
	"github.com/twstrike/coyim/Godeps/_workspace/src/github.com/twstrike/gotk3adapter/gtki"
)

type window struct {
	*bin
	internal *gtk.Window
}

type asWindow interface {
	toWindow() *window
}

func (v *window) toWindow() *window {
	return v
}

func wrapWindowSimple(v *gtk.Window) *window {
	if v == nil {
		return nil
	}
	return &window{wrapBinSimple(&v.Bin), v}
}

func wrapWindow(v *gtk.Window, e error) (*window, error) {
	return wrapWindowSimple(v), e
}

func unwrapWindow(v gtki.Window) *gtk.Window {
	if v == nil {
		return nil
	}
	return v.(asWindow).toWindow().internal
}

func (v *window) AddAccelGroup(v2 gtki.AccelGroup) {
	v.internal.AddAccelGroup(unwrapAccelGroup(v2))
}

func (v *window) GetTitle() string {
	v1, e := v.internal.GetTitle()
	if e != nil {
		return ""
	}
	return v1
}

func (v *window) IsActive() bool {
	return v.internal.IsActive()
}

func (v *window) Resize(v1, v2 int) {
	v.internal.Resize(v1, v2)
}

func (v *window) SetApplication(v2 gtki.Application) {
	v.internal.SetApplication(unwrapApplication(v2))
}

func (v *window) SetIcon(v2 gdki.Pixbuf) {
	v.internal.SetIcon(gdka.UnwrapPixbuf(v2))
}

func (v *window) SetTitle(v1 string) {
	v.internal.SetTitle(v1)
}

func (v *window) SetTitlebar(v2 gtki.Widget) {
	v.internal.SetTitlebar(unwrapWidget(v2))
}

func (v *window) SetTransientFor(v2 gtki.Window) {
	v.internal.SetTransientFor(unwrapWindow(v2))
}

func (v *window) Present() {
	v.internal.Present()
}

func (v *window) HasToplevelFocus() bool {
	return v.internal.HasToplevelFocus()
}

func (v *window) Fullscreen() {
	v.internal.Fullscreen()
}

func (v *window) Unfullscreen() {
	v.internal.Unfullscreen()
}
