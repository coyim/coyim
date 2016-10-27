package gdki

import "github.com/twstrike/gotk3adapter/glibi"

type PixbufLoader interface {
	glibi.Object

	Close() error
	GetPixbuf() (Pixbuf, error)
	Write([]byte) (int, error)
}

func AssertPixbufLoader(_ PixbufLoader) {}
