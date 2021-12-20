package gui

import (
	"os"
	"strings"
	"sync"
)

type colorSet struct {
	rosterPeerBackground                      hexColor
	rosterPeerOfflineForeground               hexColor
	rosterPeerOnlineForeground                hexColor
	rosterGroupBackground                     hexColor
	rosterAccountOnlineBackground             hexColor
	rosterAccountOfflineBackground            hexColor
	conversationOutgoingUserForeground        hexColor
	conversationIncomingUserForeground        hexColor
	conversationOutgoingTextForeground        hexColor
	conversationIncomingTextForeground        hexColor
	conversationStatusTextForeground          hexColor
	conversationOutgoingDelayedUserForeground hexColor
	conversationOutgoingDelayedTextForeground hexColor

	// these two don't seem to be used anymore
	conversationLockTypingBackground   hexColor
	conversationUnlockTypingBackground hexColor

	timestampForeground hexColor
}

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
		// - create an invisible frame and check the background and see if it is dark by default
		// - Once we have that, we should also make an icon-set, not just a color-set to keep track of all the
		// variants.
	})

	return cm.themeVariant
}

func (cm *hasColorManagement) isDarkThemeVariant() bool {
	return cm.getThemeVariant() == darkThemeVariantName
}

func (cm *hasColorManagement) currentColorSet() colorSet {
	if cm.isDarkThemeVariant() {
		return defaultDarkColorSet()
	}
	return defaultLightColorSet()
}

func defaultLightColorSet() colorSet {
	return colorSet{
		rosterPeerBackground:                      rgbFromHex("#ffffff"),
		rosterPeerOfflineForeground:               rgbFromHex("#aaaaaa"),
		rosterPeerOnlineForeground:                rgbFromHex("#000000"),
		rosterGroupBackground:                     rgbFromHex("#e9e7f3"),
		rosterAccountOnlineBackground:             rgbFromHex("#918caa"),
		rosterAccountOfflineBackground:            rgbFromHex("#d5d3de"),
		conversationOutgoingUserForeground:        rgbFromHex("#3465a4"),
		conversationIncomingUserForeground:        rgbFromHex("#a40000"),
		conversationOutgoingTextForeground:        rgbFromHex("#000000"),
		conversationIncomingTextForeground:        rgbFromHex("#000000"),
		conversationStatusTextForeground:          rgbFromHex("#aaaaaa"),
		conversationOutgoingDelayedUserForeground: rgbFromHex("#aaaaaa"),
		conversationOutgoingDelayedTextForeground: rgbFromHex("#aaaaaa"),
		conversationLockTypingBackground:          rgbFromHex("#e0e0e0"),
		conversationUnlockTypingBackground:        rgbFromHex("#f9f9f9"),
		timestampForeground:                       rgbFromHex("#aaaaaa"),
	}
}

func defaultDarkColorSet() colorSet {
	return colorSet{
		rosterPeerBackground:                      rgbFromHex("#7f7f7f"),
		rosterPeerOfflineForeground:               rgbFromHex("#aaaaaa"),
		rosterPeerOnlineForeground:                rgbFromHex("#e5e5e5"),
		rosterGroupBackground:                     rgbFromHex("#b8b6bf"),
		rosterAccountOnlineBackground:             rgbFromHex("#d5d3de"),
		rosterAccountOfflineBackground:            rgbFromHex("#918caa"),
		conversationOutgoingUserForeground:        rgbFromHex("#3465a4"),
		conversationIncomingUserForeground:        rgbFromHex("#a40000"),
		conversationOutgoingTextForeground:        rgbFromHex("#7f7f7f"),
		conversationIncomingTextForeground:        rgbFromHex("#7f7f7f"),
		conversationStatusTextForeground:          rgbFromHex("#4e9a06"),
		conversationOutgoingDelayedUserForeground: rgbFromHex("#444444"),
		conversationOutgoingDelayedTextForeground: rgbFromHex("#444444"),
		conversationLockTypingBackground:          rgbFromHex("#e0e0e0"),
		conversationUnlockTypingBackground:        rgbFromHex("#252a2c"),
		timestampForeground:                       rgbFromHex("#444444"),
	}
}
