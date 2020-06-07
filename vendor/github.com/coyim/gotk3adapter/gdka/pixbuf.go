package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/gotk3/gotk3/gdk"
)

type pixbuf struct {
	internal *gdk.Pixbuf
}

func WrapPixbufSimple(v *gdk.Pixbuf) gdki.Pixbuf {
	if v == nil {
		return nil
	}
	return &pixbuf{v}
}

func WrapPixbuf(v *gdk.Pixbuf, e error) (gdki.Pixbuf, error) {
	return WrapPixbufSimple(v), e
}

func UnwrapPixbuf(v gdki.Pixbuf) *gdk.Pixbuf {
	if v == nil {
		return nil
	}
	return v.(*pixbuf).internal
}

func (v *pixbuf) SavePNG(filename string, compression int) error {
	return v.internal.SavePNG(filename, compression)
}
