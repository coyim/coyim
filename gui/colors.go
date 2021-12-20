package gui

import (
	"os"
	"strings"
	"sync"
)

type colorSet struct {
	rosterPeerBackground                      string
	rosterPeerOfflineForeground               string
	rosterPeerOnlineForeground                string
	rosterGroupBackground                     string
	rosterAccountOnlineBackground             string
	rosterAccountOfflineBackground            string
	conversationOutgoingUserForeground        string
	conversationIncomingUserForeground        string
	conversationOutgoingTextForeground        string
	conversationIncomingTextForeground        string
	conversationStatusTextForeground          string
	conversationOutgoingDelayedUserForeground string
	conversationOutgoingDelayedTextForeground string
	conversationLockTypingBackground          string
	conversationUnlockTypingBackground        string
	timestampForeground                       string
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
		// variates.
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
		rosterPeerBackground:                      "#ffffff",
		rosterPeerOfflineForeground:               "#aaaaaa",
		rosterPeerOnlineForeground:                "#000000",
		rosterGroupBackground:                     "#e9e7f3",
		rosterAccountOnlineBackground:             "#918caa",
		rosterAccountOfflineBackground:            "#d5d3de",
		conversationOutgoingUserForeground:        "#3465a4",
		conversationIncomingUserForeground:        "#a40000",
		conversationOutgoingTextForeground:        "#000000",
		conversationIncomingTextForeground:        "#000000",
		conversationStatusTextForeground:          "#aaaaaa",
		conversationOutgoingDelayedUserForeground: "#aaaaaa",
		conversationOutgoingDelayedTextForeground: "#aaaaaa",
		conversationLockTypingBackground:          "#e0e0e0",
		conversationUnlockTypingBackground:        "#f9f9f9",
		timestampForeground:                       "#aaaaaa",
	}
}

func defaultDarkColorSet() colorSet {
	return colorSet{
		rosterPeerBackground:                      "#7f7f7f",
		rosterPeerOfflineForeground:               "#aaaaaa",
		rosterPeerOnlineForeground:                "#e5e5e5",
		rosterGroupBackground:                     "#b8b6bf",
		rosterAccountOnlineBackground:             "#d5d3de",
		rosterAccountOfflineBackground:            "#918caa",
		conversationOutgoingUserForeground:        "#3465a4",
		conversationIncomingUserForeground:        "#a40000",
		conversationOutgoingTextForeground:        "#7f7f7f",
		conversationIncomingTextForeground:        "#7f7f7f",
		conversationStatusTextForeground:          "#4e9a06",
		conversationOutgoingDelayedUserForeground: "#444444",
		conversationOutgoingDelayedTextForeground: "#444444",
		conversationLockTypingBackground:          "#e0e0e0",
		conversationUnlockTypingBackground:        "#252a2c",
		timestampForeground:                       "#444444",
	}
}
