package gui

import (
	"fmt"
	"strings"

	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
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
	if ds.fontSize == 0 {
		ds.fontSize = 12
		ds.update()
		return
	}
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
        `, ds.fontSize)
	doInUIThread(func() {
		_ = ds.provider.LoadFromData(css)
	})
}

func newDisplaySettings() *displaySettings {
	ds := &displaySettings{}
	prov, _ := g.gtk.CssProviderNew()
	ds.provider = prov
	ds.defaultFontSize = 12
	return ds
}

// TODO: currently, while using the shortcuts, the font
// changes to 6.
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
		styleContext, _ := l.GetStyleContext()
		ds := newDisplaySettings()

		styleContext.AddClass("bold-header-style")
		styleContext.AddProvider(ds.provider, 9999)

		_ = ds.provider.LoadFromData(`.bold-header-style {
			font-size: 200%;
			font-weight: 800;
		}`)
	})
}

// StyleContextable is an interface to assing css style
type StyleContextable interface {
	GetStyleContext() (gtki.StyleContext, error)
}

func providerWithCSS(s string) gtki.CssProvider {
	p, _ := g.gtk.CssProviderNew()
	_ = p.LoadFromData(s)
	return p
}

func updateWithStyle(l StyleContextable, p gtki.CssProvider) {
	sc, _ := l.GetStyleContext()
	sc.AddProvider(p, 9999)
}

type style map[string]interface{}
type styles map[string]style

func styleSelectorRules(el string, s style) string {
	return fmt.Sprintf("%s {%s}", el, inlineStyleProperties(s))
}

func providerWithStyle(el string, s style) gtki.CssProvider {
	return providerWithStyles(styles{el: s})
}

func providerWithStyles(st styles) gtki.CssProvider {
	selectors := []string{}
	for el, s := range st {
		selectors = append(selectors, styleSelectorRules(el, s))
	}
	return providerWithCSS(strings.Join(selectors, ""))
}

func inlineStyleProperties(s style) string {
	inline := []string{}
	for attr, value := range s {
		inline = append(inline, fmt.Sprintf("%s: %s;", attr, value))
	}
	return strings.Join(inline, "")
}
