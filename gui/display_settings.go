package gui

import (
	"fmt"

	"github.com/twstrike/gotk3adapter/gtki"
	"github.com/twstrike/gotk3adapter/pangoi"
)

type displaySettings struct {
	fontSize        uint
	defaultFontSize uint

	provider gtki.CssProvider
}

func (ds *displaySettings) defaultSettingsOn(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
	})
}

func (ds *displaySettings) unifiedBackgroundColor(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("currentBackgroundColor")
	})
}

func (ds *displaySettings) control(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("currentFontSetting")
	})
}

func (ds *displaySettings) shadeBackground(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider, 9999)
		styleContext.AddClass("shadedBackgroundColor")
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

.shadedBackgroundColor {
  background-color: #fafafa;
}
`, ds.defaultFontSize, ds.fontSize)
	doInUIThread(func() {
		ds.provider.LoadFromData(css)
	})
}

func newDisplaySettings() *displaySettings {
	ds := &displaySettings{}
	prov, _ := g.gtk.CssProviderNew()
	ds.provider = prov
	ds.defaultFontSize = 12
	return ds
}

func getFontSizeFrom(w gtki.Widget) uint {
	styleContext, _ := w.GetStyleContext()
	property, _ := styleContext.GetProperty2("font", gtki.STATE_FLAG_NORMAL)
	fontDescription := property.(pangoi.FontDescription)
	return uint(fontDescription.GetSize() / pangoi.PANGO_SCALE)
}

func detectCurrentDisplaySettingsFrom(w gtki.Widget) *displaySettings {
	ds := newDisplaySettings()
	ds.fontSize = getFontSizeFrom(w)
	return ds
}

func addBoldHeaderStyle(l gtki.Label) {
	doInUIThread(func() {
		c, _ := l.GetStyleContext()
		ds := newDisplaySettings()
		c.AddClass("bold-header-style")
		c.AddProvider(ds.provider, 9999)
		ds.provider.LoadFromData(`.bold-header-style {
			font-size: 200%;
			font-weight: 800;
		}`)
	})
}

func addStyle(w gtki.Widget, className string, css string) {
	doInUIThread(func() {
		c, _ := w.GetStyleContext()
		ds := newDisplaySettings()
		c.AddProvider(ds.provider, 9999)
		c.AddClass(className)
		ds.provider.LoadFromData(css)
	})
}
