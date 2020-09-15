package gui

import (
	"sync"

	"github.com/coyim/gotk3adapter/gtki"

	"github.com/coyim/gotk3adapter/glibi"
)

const (
	indexPixbufColumn = iota
)

var pixbufType func() glibi.Type = func() func() glibi.Type {
	var onlyOnce sync.Once
	var tp glibi.Type

	readPixbufType := func() {
		builder := newBuilderFromString(`
<interface>
	<object id="storeOfColumnTypes" class="GtkListStore">
		<columns>
			<column type="GdkPixbuf"/>
		</columns>
	</object>
</interface>
`)
		store := builder.get("storeOfColumnTypes").(gtki.ListStore)
		tp = store.GetColumnType(indexPixbufColumn)
	}

	return func() glibi.Type {
		onlyOnce.Do(readPixbufType)
		return tp
	}
}()
