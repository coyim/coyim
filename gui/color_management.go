package gui

import (
	"os"
	"strings"
	"sync"

	"github.com/coyim/gotk3adapter/gtki"
)

type hasColorManagement struct {
	themeVariant          string
	calculateThemeVariant sync.Once
}

const darkThemeVariantName = "dark"
const lightThemeVariantName = "light"

func (cm *hasColorManagement) getThemeVariant() string {
	cm.calculateThemeVariant.Do(func() {
		cm.themeVariant = lightThemeVariantName
		gtkTheme := os.Getenv("GTK_THEME")
		if gtkTheme != "" {
			toks := strings.Split(gtkTheme, ":")
			variant := toks[len(toks)-1:][0]
			if variant == darkThemeVariantName {
				cm.themeVariant = variant
				return
			}
		}

		settings, err := g.gtk.SettingsGetDefault()
		if err != nil {
			panic(err)
		}

		prefDark, _ := settings.GetProperty("gtk-application-prefer-dark-theme")
		if val, ok := prefDark.(bool); val && ok {
			cm.themeVariant = darkThemeVariantName
			return
		}

		// TODO: we should do two things here
		// - check the current theme name, and see if it ends with -dark or _dark - not just splitting on the ":" as above

		bgcd := newBackgroundColorDetectionInvisibleListBox()
		styleContext, _ := bgcd.lb.GetStyleContext()
		bc, _ := styleContext.GetProperty2("background-color", gtki.STATE_FLAG_NORMAL)
		bgcd.lb.Destroy()
		if rgbFromGetters(bc.(rgbaGetters)).isDark() {
			cm.themeVariant = darkThemeVariantName
		}

	})

	return cm.themeVariant
}

func (cm *hasColorManagement) isDarkThemeVariant() bool {
	return cm.getThemeVariant() == darkThemeVariantName
}
