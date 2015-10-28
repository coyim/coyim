package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

type displaySettings struct {
	fontSize uint

	provider *gtk.CssProvider
}

func (ds *displaySettings) unifiedBackgroundColor(w *gtk.Widget) {
	glib.IdleAdd(func() bool {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("currentBackgroundColor")
		return false
	})
}

func (ds *displaySettings) control(w *gtk.Widget) {
	glib.IdleAdd(func() bool {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("currentFontSetting")
		return false
	})
}

func (ds *displaySettings) increaseFontSize() {
	ds.fontSize++
	ds.update()
}

func (ds *displaySettings) decreaseFontSize() {
	ds.fontSize--
	ds.update()
}

func (ds *displaySettings) update() {
	css := fmt.Sprintf(`
.currentFontSetting {
  font-size: %dpx;
}

.currentBackgroundColor {
  background-color: #fff;
}
`, ds.fontSize)
	glib.IdleAdd(func() bool {
		ds.provider.LoadFromData(css)
		return false
	})
}

func newDisplaySettings() *displaySettings {
	ds := &displaySettings{}
	prov, _ := gtk.CssProviderNew()
	ds.provider = prov
	return ds
}

func detectCurrentDisplaySettingsFrom(w *gtk.Widget) *displaySettings {
	styleContext, _ := w.GetStyleContext()
	property, _ := styleContext.GetProperty("font", gtk.STATE_FLAG_NORMAL)
	fontDescription := property.(*pango.FontDescription)

	size := uint(fontDescription.GetSize() / pango.PANGO_SCALE)
	ds := newDisplaySettings()
	ds.fontSize = size
	return ds
}
