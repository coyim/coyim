package gdki

import "github.com/coyim/gotk3adapter/glibi"

type PixbufLoader interface {
	glibi.Object

	Close() error
	GetPixbuf() (Pixbuf, error)
	SetSize(int, int)
	Write([]byte) (int, error)
}

func AssertPixbufLoader(_ PixbufLoader) {}
