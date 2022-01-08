package gui

import (
	"fmt"
	"strings"

	"github.com/coyim/coyim/gui/css"
	"github.com/coyim/gotk3adapter/gtki"
	"github.com/coyim/gotk3adapter/pangoi"
	"github.com/gotk3/gotk3/gtk"
)

type displaySettings struct {
	fontSize        uint
	defaultFontSize uint

	provider             *cssProvider
	globalColorsProvider *cssProvider
}

func (ds *displaySettings) defaultSettingsOn(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider.provider, 9999)
	})
}

func (ds *displaySettings) control(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider.provider, 9999)
		styleContext.AddClass("currentFontSetting")
	})
}

func (ds *displaySettings) shadeBackground(w gtki.Widget) {
	doInUIThread(func() {
		styleContext, _ := w.GetStyleContext()
		styleContext.AddProvider(ds.provider.provider, 9999)
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
		ds.provider.load("current font", css)
	})
}

type cssProvider struct {
	provider gtki.CssProvider
	l        withLog
}

func (c *cssProvider) load(name, s string) {
	e := c.provider.LoadFromData(s)
	if e != nil {
		c.l.Log().WithError(e).WithField("name", name).Error("couldn't load CSS data")
	}
}

func newCSSProvider(hl withLog) *cssProvider {
	prov, _ := g.gtk.CssProviderNew()
	return &cssProvider{
		provider: prov,
		l:        hl,
	}
}

func newDisplaySettings(hl withLog) *displaySettings {
	ds := &displaySettings{}

	ds.provider = newCSSProvider(hl)
	ds.defaultFontSize = 12

	ds.globalColorsProvider = newCSSProvider(hl)

	doInUIThread(func() {
		// TODO: we need to update loading of color definitions
		// based on dark or light theme
		ds.globalColorsProvider.load("color definitions", css.Get("light/colors.css"))
		addGlobalProvider(ds.globalColorsProvider.provider)
	})

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

func detectCurrentDisplaySettingsFrom(hl withLog, w gtki.Widget) *displaySettings {
	ds := newDisplaySettings(hl)
	ds.fontSize = getFontSizeFrom(w)
	return ds
}

func addBoldHeaderStyle(hl withLog, l gtki.Label) {
	doInUIThread(func() {
		styleContext, _ := l.GetStyleContext()
		ds := newDisplaySettings(hl)

		styleContext.AddClass("bold-header-style")
		styleContext.AddProvider(ds.provider.provider, 9999)

		ds.provider.load("bold header", css.Get("bold_header_style.css"))
	})
}

type styleContextable interface {
	GetStyleContext() (gtki.StyleContext, error)
}

func providerFromCSSFile(wl withLog, msg, file string) gtki.CssProvider {
	return providerWithCSS(wl, msg, css.Get(file))
}

func providerWithCSS(wl withLog, msg, s string) gtki.CssProvider {
	p := newCSSProvider(wl)
	p.load(msg, s)
	return p.provider
}

const styleProviderHighPriority = gtk.STYLE_PROVIDER_PRIORITY_USER * 10

func addGlobalProvider(p gtki.CssProvider) {
	screen, e := g.gdk.ScreenGetDefault()
	panicOnDevError(e)
	g.gtk.AddProviderForScreen(screen, p, styleProviderHighPriority)
}

func updateWithStyle(l styleContextable, p gtki.CssProvider) {
	if sc, err := l.GetStyleContext(); err == nil {
		sc.AddProvider(p, styleProviderHighPriority)
	}
}

func updateWithStyles(l styleContextable, p gtki.CssProvider) {
	sc, _ := l.GetStyleContext()
	screen, _ := sc.GetScreen()
	g.gtk.AddProviderForScreen(screen, p, styleProviderHighPriority)
}

type style map[string]interface{}
type styles map[string]style

type nestedStyles struct {
	rootSelector string
	rootStyle    style
	nestedStyles styles
}

func (nst *nestedStyles) toStyles() styles {
	selector := nst.rootSelector

	ret := styles{
		selector: nst.rootStyle,
	}

	for s, n := range nst.nestedStyles {
		ret[nestedCSSRules(selector, s)] = n
	}

	return ret
}

func styleSelectorRules(el string, s style) string {
	return fmt.Sprintf("%s {%s}", el, inlineStyleProperties(s))
}

func providerWithStyle(wl withLog, name string, el string, s style) gtki.CssProvider {
	return providerWithStyles(wl, name, styles{el: s})
}

func providerWithStyles(wl withLog, name string, st styles) gtki.CssProvider {
	selectors := []string{}
	for el, s := range st {
		selectors = append(selectors, styleSelectorRules(el, s))
	}
	return providerWithCSS(wl, name, strings.Join(selectors, ""))
}

func inlineStyleProperties(s style) string {
	inline := []string{}
	for attr, value := range s {
		inline = append(inline, fmt.Sprintf("%s: %s;", attr, value))
	}
	return strings.Join(inline, "")
}
