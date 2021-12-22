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

func doesThemeNameIndicateDarkness(name string) bool {
	parts := strings.Split(name, ":")
	if len(parts) < 2 {
		return false
	}
	variant := parts[len(parts)-1]
	return variant == darkThemeVariantName
}

func (cm *hasColorManagement) detectDarkThemeFromEnvironmentVariable() bool {
	gtkTheme := os.Getenv("GTK_THEME")
	return doesThemeNameIndicateDarkness(gtkTheme)
}

func (cm *hasColorManagement) detectDarkThemeFromGTKSettings() bool {
	// TODO: this might not be safe to do outside the UI thread
	settings, err := g.gtk.SettingsGetDefault()
	if err != nil {
		panic(err)
	}

	prefDark, _ := settings.GetProperty("gtk-application-prefer-dark-theme")
	val, ok := prefDark.(bool)
	return val && ok
}

func (cm *hasColorManagement) detectDarkThemeFromGTKListBoxBackground() bool {
	// TODO: this is not safe to do outside the UI thread
	bgcd := newBackgroundColorDetectionInvisibleListBox()
	styleContext, _ := bgcd.lb.GetStyleContext()
	bc, _ := styleContext.GetProperty2("background-color", gtki.STATE_FLAG_NORMAL)
	bgcd.lb.Destroy()
	return rgbFromGetters(bc.(rgbaGetters)).isDark()
}

func (cm *hasColorManagement) isDarkTheme() bool {
	return cm.detectDarkThemeFromEnvironmentVariable() ||
		cm.detectDarkThemeFromGTKSettings() ||
		cm.detectDarkThemeFromGTKListBoxBackground()
}

func (cm *hasColorManagement) actuallyCalculateThemeVariant() {
	if cm.isDarkTheme() {
		cm.themeVariant = darkThemeVariantName
	} else {
		cm.themeVariant = lightThemeVariantName
	}

	// - check the current theme name, and see if it ends with -dark or _dark - not just splitting on the ":" as above
}

func (cm *hasColorManagement) getThemeVariant() string {
	cm.calculateThemeVariant.Do(cm.actuallyCalculateThemeVariant)
	return cm.themeVariant
}

func (cm *hasColorManagement) isDarkThemeVariant() bool {
	return cm.getThemeVariant() == darkThemeVariantName
}
