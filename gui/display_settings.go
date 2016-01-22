package gui

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/gotk3/pango"
)

type displaySettings struct {
	fontSize        uint
	defaultFontSize uint

	provider *gtk.CssProvider
}

func (ds *displaySettings) defaultSettingsOn(w *gtk.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
	})
}

func (ds *displaySettings) unifiedBackgroundColor(w *gtk.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("currentBackgroundColor")
	})
}

func (ds *displaySettings) control(w *gtk.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("currentFontSetting")
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
* {
  -GtkCheckMenuItem-indicator-size: 16;
}

.currentFontSetting {
  font-size: %dpx;
}

.currentBackgroundColor {
  background-color: #fff;
}
`, ds.defaultFontSize, ds.fontSize)
	doInUIThread(func() {
		ds.provider.LoadFromData(css)
	})
}

func newDisplaySettings() *displaySettings {
	ds := &displaySettings{}
	prov, _ := gtk.CssProviderNew()
	ds.provider = prov
	ds.defaultFontSize = 12
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
