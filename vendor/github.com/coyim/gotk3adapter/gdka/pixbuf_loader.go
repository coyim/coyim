package gdka

import (
	"github.com/coyim/gotk3adapter/gdki"
	"github.com/coyim/gotk3adapter/gliba"
	"github.com/gotk3/gotk3/gdk"
)

type pixbufLoader struct {
	*gliba.Object
	internal *gdk.PixbufLoader
}

func WrapPixbufLoaderSimple(v *gdk.PixbufLoader) gdki.PixbufLoader {
	if v == nil {
		return nil
	}
	return &pixbufLoader{gliba.WrapObjectSimple(v.Object), v}
}

func WrapPixbufLoader(v *gdk.PixbufLoader, e error) (gdki.PixbufLoader, error) {
	return WrapPixbufLoaderSimple(v), e
}

func UnwrapPixbufLoader(v gdki.PixbufLoader) *gdk.PixbufLoader {
	if v == nil {
		return nil
	}
	return v.(*pixbufLoader).internal
}

func (v *pixbufLoader) Close() error {
	return v.internal.Close()
}

func (v *pixbufLoader) GetPixbuf() (gdki.Pixbuf, error) {
	return WrapPixbuf(v.internal.GetPixbuf())
}

func (v *pixbufLoader) SetSize(width, height int) {
	v.internal.SetSize(width, height)
}

func (v *pixbufLoader) Write(b []byte) (int, error) {
	return v.internal.Write(b)
}
