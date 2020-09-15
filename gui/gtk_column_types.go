package gui

import (
	"sync"

	"github.com/coyim/gotk3adapter/glibi"

	"github.com/coyim/gotk3adapter/gtki"
)

const (
	indexPixbufColumn = iota
)

type gtkColumnTypes struct {
	store gtki.ListStore `gtk-widget:"storeOfColumnTypes"`
}

var pixbufType func() glibi.Type = func() func() glibi.Type {
	var onlyOnce sync.Once
	var tp glibi.Type

	readPixbufType := func() {
		ct := &gtkColumnTypes{}
		builder := newBuilder("GTKColumnTypes")
		panicOnDevError(builder.bindObjects(ct))
		tp = ct.store.GetColumnType(indexPixbufColumn)
	}

	return func() glibi.Type {
		onlyOnce.Do(readPixbufType)
		return tp
	}
}()
