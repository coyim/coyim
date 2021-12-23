package gui

import (
	"os"
	"strings"
	"sync"

	"github.com/coyim/gotk3adapter/glibi"
	"github.com/coyim/gotk3adapter/gtki"
)

type hasColorManagement struct {
	themeVariant          string
	calculateThemeVariant sync.Once
}

const darkThemeVariantName = "dark"
const lightThemeVariantName = "light"

func isDarkVariantNameBasedOnSeparator(name, sep string) bool {
	parts := strings.Split(name, sep)
	if len(parts) < 2 {
		return false
	}
	variant := parts[len(parts)-1]
	return variant == darkThemeVariantName
}

func doesThemeNameIndicateDarkness(name string) bool {
	return isDarkVariantNameBasedOnSeparator(name, ":") ||
		isDarkVariantNameBasedOnSeparator(name, "-") ||
		isDarkVariantNameBasedOnSeparator(name, "_")
}

func (cm *hasColorManagement) detectDarkThemeFromEnvironmentVariable() bool {
	gtkTheme := os.Getenv("GTK_THEME")
	return doesThemeNameIndicateDarkness(gtkTheme)
}

func (cm *hasColorManagement) getGTKSettings() gtki.Settings {
	settings, err := g.gtk.SettingsGetDefault()
	if err != nil {
		panic(err)
	}
	return settings
}

func (cm *hasColorManagement) getGSettings() glibi.Settings {
	return g.glib.SettingsNew("org.gnome.desktop.interface")
}

func (cm *hasColorManagement) getThemeNameFromGTKSettings() string {
	// TODO: this might not be safe to do outside the UI thread
	themeName, _ := cm.getGTKSettings().GetProperty("gtk-theme-name")
	val, _ := themeName.(string)
	return val
}

func (cm *hasColorManagement) getThemeNameFromGSettings() string {
	// TODO: this might not be safe to do outside the UI thread
	return cm.getGSettings().GetString("gtk-theme")
}

func (cm *hasColorManagement) detectDarkThemeFromGTKSettings() bool {
	// TODO: this might not be safe to do outside the UI thread
	prefDark, _ := cm.getGTKSettings().GetProperty("gtk-application-prefer-dark-theme")
	val, ok := prefDark.(bool)
	return val && ok
}

func (cm *hasColorManagement) detectDarkThemeFromGTKSettingsThemeName() bool {
	return doesThemeNameIndicateDarkness(cm.getThemeNameFromGTKSettings())
}

func (cm *hasColorManagement) detectDarkThemeFromGSettingsThemeName() bool {
	return doesThemeNameIndicateDarkness(cm.getThemeNameFromGSettings())
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
		cm.detectDarkThemeFromGTKSettingsThemeName() ||
		cm.detectDarkThemeFromGSettingsThemeName() ||
		cm.detectDarkThemeFromGTKListBoxBackground()
}

func (cm *hasColorManagement) actuallyCalculateThemeVariant() {
	if cm.isDarkTheme() {
		cm.themeVariant = darkThemeVariantName
	} else {
		cm.themeVariant = lightThemeVariantName
	}
}

func (cm *hasColorManagement) getThemeVariant() string {
	cm.calculateThemeVariant.Do(cm.actuallyCalculateThemeVariant)
	return cm.themeVariant
}

func (cm *hasColorManagement) isDarkThemeVariant() bool {
	return cm.getThemeVariant() == darkThemeVariantName
}
