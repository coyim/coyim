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

func wrapPixbufLoaderSimple(v *gdk.PixbufLoader) *pixbufLoader {
	if v == nil {
		return nil
	}
	return &pixbufLoader{gliba.WrapObjectSimple(v.Object), v}
}

func wrapPixbufLoader(v *gdk.PixbufLoader, e error) (*pixbufLoader, error) {
	return wrapPixbufLoaderSimple(v), e
}

func unwrapPixbufLoader(v gdki.PixbufLoader) *gdk.PixbufLoader {
	if v == nil {
		return nil
	}
	return v.(*pixbufLoader).internal
}

func (v *pixbufLoader) Close() error {
	return v.internal.Close()
}

func (v *pixbufLoader) GetPixbuf() (gdki.Pixbuf, error) {
	return wrapPixbuf(v.internal.GetPixbuf())
}

func (v *pixbufLoader) SetSize(width, height int) {
	v.internal.SetSize(width, height)
}

func (v *pixbufLoader) Write(b []byte) (int, error) {
	return v.internal.Write(b)
}
